package gtfsparser

// The option mechanism is used to set flags/values which modify the parsing behaviour. For instance, to ignore certain parsing errors.
// This is necessary as not all organisations choose to implement GTFS correctly, and some mandate special behaviour.

type Opt struct {
  enabled bool
  value   string
}

type OptType int

const (
  OPT_IGNORE_SHAPE_TRIP_MAPPING_ERROR OptType = iota
)
