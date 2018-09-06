LDFLAGS = -extldflags -static -s -w

build:
	mkdir -p ./buildtarget
	env CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o ./buildtarget/xlscalendar ./cmd/xlscalendar

clean:
	rm -rf ./buildtarget/*
