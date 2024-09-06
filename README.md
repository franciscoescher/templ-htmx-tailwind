# TODO Templ + Htmx + Tailwind CSS

This is a simple TODO project to run [templ](https://github.com/a-h/templ) + [htmx](https://htmx.org/) + [tailwind css](https://tailwindcss.com/).

It uses [yarn](https://yarnpkg.com/) + [vite](https://vitejs.dev/) to build tailwind locally and allow setup of local configs.

It uses [air](https://github.com/air-verse/air) to auto reload files.

Install yarn+air prior to running the project.

## Run with hot reload

`air`

## Run manually

Build js files:

`yarn build`

Generate templ code:

`templ generate`

Run the server:

`go run .`

Or all at once:

`yarn build && templ generate && go run .`
