package fmt

// Fmt prints a message to stdout or stderr.
type Fmt interface {
	Print(i ...interface{})
	Println(i ...interface{})
	Printf(format string, i ...interface{})
	PrintErr(i ...interface{})
	PrintErrln(i ...interface{})
	PrintErrf(format string, i ...interface{})
}

// DataWriter is the interface that wraps the basic WriteData method.
type DataWriter interface {
	WriteData(v interface{}) error
}

// DataWriterProvider provides DataWriter.
type DataWriterProvider interface {
	DataWriter() DataWriter
}
