go-license
==========

A license management utility for programs written in Golang.

This program handles identifying software licenses and standardizing on a short,
abbreviated name for each known license type.

## Enforcement

License identifier enforcement is not strict. This makes it possible to warn
when an unrecognized license type is used, encouraging either conformance or an
update to the list of known licenses. There is no way we can know all types of
licenses.

## License guessing

This program also provides naive license guessing based on the license body
(text). This makes it easy to just throw a blob of text in and get a
standardized license identifier string out.
