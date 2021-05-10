# omdb
CLI tool to query https://www.omdbapi.com/

## add to .bashrc
```export OMDBAPIKEY=omdbkey```

## usage
```omdb search terminator
omdb search -type movie terminator
omdb search -type series terminator
omdb title -t terminator
omdb title terminator
omdb title -type series terminator
omdb title -id tt0851851
```

the template used to generate films.html file when running **omdb search moviename** comes from https://codepen.io/KungFuSucker/pen/xaBvep