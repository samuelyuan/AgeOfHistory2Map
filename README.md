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

* The scenario file consists of a list of civilization tags, which are represented as strings. The file is located in the scenario folder and has the same name as the scenario without any extensions.
* The province owners file consists of a list of province owners, which are represented as integers. Each province's owner corresponds to a civilization. The file is located in the scenario folder and ends in "_PD".
* The civilizations file contains a list of civilizations with the tag and the color in RGB format. This is the color that will show on the map for all provinces owned by a civilization. The files are located in "game/civilizations".
* The provinces file contains the geometry of the province. It has a separate list of X and Y coordinates which can be used to draw the shape. The files are located in "map/Earth/data/provinces".
