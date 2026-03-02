package version

// Override via -ldflags:
// -X 'github.com/revolver-sys/lanstreamd/internal/version.Version=v0.0.1'
// -X 'github.com/revolver-sys/lanstreamd/internal/version.Commit=abcdef'
var (
	Version = "dev"
	Commit  = "none"
)
