# Chaos Gate: Daemonhunters — Save File Format

Reverse-engineered notes for `*.gksave` files (Steam/Proton appID `1611910`,
Unity game by Complex Games).

## 1. On-disk layout

A save is three `\r\n`-separated chunks:

```
<header JSON> \r\n <state> \r\n <combatState> \r\n
```

| Chunk | Contents |
|-------|----------|
| `chunks[0]` | **Header** — plain UTF-8 JSON (not obfuscated). See `internal.Header`. |
| `chunks[1]` | **State** — obfuscated JSON. First byte is always `194` (`0xC2`). |
| `chunks[2]` | **Combat state** — present (non-empty) only for in-combat saves. |

The trailing `\r\n` produces an empty 4th chunk on split; the loader reads only
the first three (`internal/files/manager.go: Load`).

> ⚠️ The split assumes the combat chunk never itself contains `\r\n`. If it can,
> `bytes.Split` would truncate `chunks[2]`. Not observed in any sample, but
> unverified.

## 2. State obfuscation (`encodeDecode`)

The state chunk is UTF-8 text where each rune's codepoint has had its nibbles
swapped (`(v&0x0F0F0F0F)<<4 | (v&0xF0F0F0F0)>>4`), skipping swaps that would
produce an invalid rune. The transform is its own inverse, so the same function
both decodes and re-encodes (`internal/files/encoder.go`). After decoding,
`chunks[1]` is normal JSON.

## 3. State structure

```jsonc
{
  "topRecord": <LinearRecord>,        // GreyKnights.PlayerState
  "linearInstanceIds": [int, ...],    // parallel to linearRecords
  "linearRecords": [<LinearRecord>, ...]
}
```

`linearInstanceIds[i]` is the in-game object id for `linearRecords[i]`
(negative numbers; new ids are minted as `min(existing) - 1`).

### LinearRecord — double-encoded objects

```jsonc
{
  "typeName": "GreyKnights.KnightState",
  "assetName": "",
  "serializedContents": "{\"givenName\":\"...\", ...}"   // a JSON *string* whose
}                                                         // contents are JSON
```

`serializedContents` is a JSON **string** containing the object's JSON. The
marshaler (`internal/marshaler.go`) decodes it on load and re-encodes on save.

Only the ~17 type names in `typeNameToObject` are parsed into Go structs
(`internal/objects/`); all other records are kept as raw bytes and round-trip
untouched. **Within modeled structs, fields the tool doesn't use are typed
`json.RawMessage`, so they are preserved verbatim**

## 4. Record-type catalog

Strategic map (the bulk): `StarMapLinkModel` (7276), `StarMapNodeModel` (3468),
`StarMapLocation`, `StarMapRoute`, `StarMapMission`, `StarMapWarpStormModel`,
`StarMapManagerSaveState`, `StarMapModel`, `StarMapPlayerShip`,
`StarMapEnemyShip`, `StarMapMissionSaveState`.

Narrative: `ArticyVariableSaveState` (4760 — story flag variables),
`TimelineEventOccasion`, `TimelineEventHolder`, plus many `*Consequence` /
`*Occasion` / `*Outcome` rule records.

Ship / campaign singletons (1 each per save): `GrandmasterSaveState`,
`GameUnlocksSaveState`, `KnightsSaveState`, `ArmourySaveState`,
`CurrencySaveState`, `TimeManagerSaveState`, `ShipStatusManagerSaveState`,
`ShipResearchSaveState`, `ShipConstructionSaveState`, `HangarSaveState`, etc.

Units / combat: `KnightState`, `DreadnoughtState`, `*AssassinState`,
`EnemyState`, `MissionState`.

## 5. Ship vs battle saves

| | Ship save | Battle save |
|--|-----------|-------------|
| `header.location` | `"COMMON_Baleful_Edict"` | battle map, e.g. `"CorruptedVessel_Export"` |
| `header.saveName` | `"On Ship - N"` | `"In Combat - N"` |
| combat chunk | empty | **present** |
| strategic state | full | **also full** |

**Key insight:** a battle save still contains the *entire* strategic state. The
active mission is a normal `StarMapMission` (its `mapName` matches the header
`location`, minus the `_Export` scene suffix).

## 6. Campaign-flow fields (feature targets)

- **`LoseGameOccasion`** — a campaign loss condition: `occasionKey:<key>`,
  `triggerTime:<day>` (saw 368/409/487/554/600). When the day reaches
  `triggerTime` you auto-lose via the GameOver cutscene.
- **`StarMapMissionSaveState`** — `bloomEruptionNumber`,
  `daysRemainingToNextBloomEruption`, `weekThatKoramarWasDefeated`,
  `act3StartDay`, `dateThatThreatLevelStartedIncreasing`.
- **`GrandmasterSaveState`** — `daysToNextReport`, `reportNumber`,
  `reportNumberInCurrentAct`.
- **`GameUnlocksSaveState.unlocks[].id`** — story/progress flags, e.g.
  `Koramar_Mission_Defeated`, `Poxus_Undefeated`, `Necrosus_Undefeated`,
  `CroweAvailable`, `Assassins_Unlocked`, `Purity_Seals_Unlocked`.

## 7. Gotchas

- **Quoting:** the marshaler uses `strconv.Quote`/`Unquote` for the
  string layer. These are Go-syntax, not JSON-syntax (e.g. JSON `\/` is not
  valid Go), but every valid sample round-trips because the inner payload is
  already-escaped ASCII JSON. `Unquote`'s error is ignored; corrupt content
  still surfaces later as `ErrWrongSaveFileFormat`. A `json.Marshal`/`Unmarshal`
  string layer would be the JSON-correct primitive if hardening is wanted.
- Numbers re-marshal without trailing `.0` (`600.0` → `600`); the game's parser
  is tolerant.

## 8. How to inspect saves

Decoding needs the unexported `encodeDecode`, so write a `package files` test:

```go
data, _ := os.ReadFile(path)
chunks := bytes.Split(data, []byte("\r\n"))
decoded := encodeDecode(chunks[1])   // chunks[1][0] == 194
// decoded is normal JSON; unmarshal into internal.State, or parse loosely and
// strconv.Unquote each record's serializedContents to read the object JSON.
```
