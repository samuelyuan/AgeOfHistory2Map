# AgeOfHistory2Map

This program will render various maps found in Age of History 2, formerly known as Age of Civilizations 2. You must have a copy of the original game data stored in the data/ folder in order to generate the map images.

## Command-Line Usage

```
./AgeOfHistory2Map.exe -map=[map name (default is regions)]
```

Example
```
./AgeOfHistory2Map.exe -map=modernworld
```

## Examples

<div style="display:inline-block;">
<img src="https://raw.githubusercontent.com/samuelyuan/AgeOfHistory2Map/master/screenshots/modernworld.png" alt="modernworld" width="510" height="300" />
<img src="https://raw.githubusercontent.com/samuelyuan/AgeOfHistory2Map/master/screenshots/1440.png" alt="1440" width="510" height="300" />
<img src="https://raw.githubusercontent.com/samuelyuan/AgeOfHistory2Map/master/screenshots/regions.png" alt="regions" width="510" height="300" />
</div>


## File Format

The map data from Age of History 2 is split into multiple files. This program assumes the user is trying to load a scenario on the Earth map. Each scenario is stored under the folder with the scenario name in "map/Earth/scenarios/". Each file is a serialized java object.

| File Type | Description | Data Format | Location |
|-----------|-------------|-------------|----------|
| Scenario file | List of civilization tags | Strings | Scenario folder, same name as scenario (no extension) |
| Province owners file | List of province owners | Integers (each corresponds to a civilization) | Scenario folder, ends in "_PD" |
| Civilizations file | List of civilizations with tag and color | Tag and RGB color format | "game/civilizations" |
| Provinces file | Geometry of provinces | Separate lists of X and Y coordinates | "map/Earth/data/provinces" |
