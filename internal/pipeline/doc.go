// Package pipeline provides a high-level orchestration layer for logslice.
//
// It composes the lower-level packages — lineparser, filter, slicer, output,
// and stats — into a single Pipeline type that can be driven from a CLI
// command or any other caller with a config.Config value.
//
// Typical usage:
//
//	cfg, err := config.FromFlags(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	pl, err := pipeline.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	snap, err := pl.Run()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("wrote %d lines\n", snap.Written)
package pipeline
