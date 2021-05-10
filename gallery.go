package main

const templ = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="ie=edge">
<title>Movies</title>
<style>
@import url('https://fonts.googleapis.com/css?family=Roboto');

* {
    box-sizing: border-box;
}

body {
    /* background: aqua; */
    background: #111116;
    font-family: 'Roboto';
}

body,
h1 {
    margin: 0;
}

header>h1 {
    text-align: center;
    margin: 10px 0;
    color: #18FF92;
}

.container {
    color: white;
    display: flex;
    flex-wrap: wrap;
    /* justify-content: space-between */
  overflow:hidden;
}

.movie-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    /* border: 1px solid #18FF92; */
    margin: 20px;
    padding-top: 10px;
    /* flex: 0 0 calc(25% - 40px); */
    flex: 0 0 calc(25% - 40px);
    /* flex: 0 0 25%; */
    /* margin-right:auto;
    margin-right:left; */
/*     overflow:hidden; */
}


.movie-item img {
    width: 100%;
    transition: transform .5s
}

img:hover {
    transform: scale(1.2);
    /* z-index:9999; */
    /* (150% zoom - Note: if the zoom is too large, it will go outside of the viewport) */
}



@media (max-width:1125px) {
    .movie-item {
        flex: 0 0 calc(33.33333% - 40px);
    }
}

@media (max-width:856px) {
    .container {
        justify-content: center;
    }
    .movie-item {
        /* margin-left:auto; */
        /* margin-right:auto; */
        flex: 0 0 calc(50% - 40px);
    }
}
</style>
</head>
<body>
<header>
<h1>Movies</h1>
</header>
<div class="container">
{{range .Search}}
<div class="movie-item">
<a href="{{ .IMDBID | printf "%s%s" "https://www.imdb.com/title/"}}">
{{if .Poster | hasPoster}}
<img src="{{.Poster}}" alt="{{.Title}}">
{{else}}
<img src="https://upload.wikimedia.org/wikipedia/commons/6/64/Poster_not_available.jpg" alt="{{.Title}}">
{{end}}
</a>
<p>{{.Title}} {{.Year}}</p>
</div>
{{end}}
</div>
</body>
</html>`

func hasPoster(s string) bool {
	return s != "N/A"
}
