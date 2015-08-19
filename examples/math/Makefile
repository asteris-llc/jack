.PHONY: clean

all: math math-add math-subtract math-multiply math-divide

math: math.go
	go build math.go

math-add: math-add.go
	go build math-add.go

math-subtract: math-subtract.go
	go build math-subtract.go

math-multiply: math-multiply.go
	go build math-multiply.go

math-divide: math-divide.go
	go build math-divide.go

clean:
	rm math math-add math-subtract math-multiply math-divide || true
