
TAILWIND_INPUT=./servers/web/styles/tailwind.css
TAILWIND_OUTPUT=./servers/web/static/main.css

run-local:
	docker-compose up --build

templ:
	cd ./servers/web && go tool templ generate

tailwind:
	npx @tailwindcss/cli -i $(TAILWIND_INPUT) -o $(TAILWIND_OUTPUT)