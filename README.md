# GOFIXEDFIELD

[![Build Status](https://secure.travis-ci.org/jbuchbinder/gofixedfield.png)](http://travis-ci.org/jbuchbinder/gofixedfield)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbuchbinder/gofixedfield)](https://goreportcard.com/report/github.com/jbuchbinder/gofixedfield)
[![GoDoc](https://godoc.org/github.com/jbuchbinder/gofixedfield?status.png)](https://godoc.org/github.com/jbuchbinder/gofixedfield)

Go library to deal with extracting fixed field form values using
struct tags.

##European-styled numbers
To parse documents that use a comma "," instead of the decimal point, just set to true the corresponding global variable:
`gofixedfield.DecimalComma = true`
