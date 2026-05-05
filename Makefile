.PHONY: setup test check validate scrape backfill charts dashboard dashboard-build open-dashboard site deploy-pages all

DASHBOARD_FILE ?= output/charts/giessen_dashboard.html
SITE_CNAME ?=

setup:
	go mod download

test:
	go test ./...

check:
	go test ./...
	go run ./cmd/validate

validate:
	go run ./cmd/validate $(ARGS)

scrape:
	go run ./cmd/scrape $(ARGS)

backfill:
	go run ./cmd/backfill $(ARGS)

charts:
	go run ./cmd/charts $(ARGS)

dashboard: dashboard-build open-dashboard

dashboard-build:
	go run ./cmd/dashboard $(ARGS)

open-dashboard:
	@file="$$(pwd)/$(DASHBOARD_FILE)"; \
	if command -v open >/dev/null 2>&1; then \
		open "$$file"; \
	elif command -v xdg-open >/dev/null 2>&1; then \
		xdg-open "$$file" >/dev/null 2>&1 & \
	else \
		echo "Dashboard geschrieben: $$file"; \
	fi

site: charts dashboard-build
	rm -rf public
	mkdir -p public/charts public/data
	touch public/.nojekyll
	if [ -n "$(SITE_CNAME)" ]; then echo "$(SITE_CNAME)" > public/CNAME; fi
	cp $(DASHBOARD_FILE) public/index.html
	cp output/charts/* public/charts/
	cp output/metadata.json output/places.csv public/data/

deploy-pages: site
	@tmp=$$(mktemp -d); \
	remote=$${DEPLOY_REMOTE:-$$(git remote get-url origin)}; \
	git clone --quiet --branch gh-pages --single-branch $$remote $$tmp; \
	git -C $$tmp rm -r --ignore-unmatch . >/dev/null; \
	cp -R public/. $$tmp/; \
	git -C $$tmp add -A; \
	if git -C $$tmp diff --cached --quiet; then \
		echo "gh-pages ist bereits aktuell"; \
	else \
		git -C $$tmp commit -m "Deploy GitHub Pages site"; \
		git -C $$tmp push origin gh-pages; \
	fi; \
	rm -rf $$tmp

all:
	go run ./cmd/scrape --postcodes all
	go run ./cmd/charts --png
	go run ./cmd/dashboard
