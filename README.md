# System Health Checker

A lightweight, modular Go CLI utility for monitoring system health in real-time. It collects essential system metrics (CPU, memory, disk usage, and optional process monitoring) and displays them through human-friendly table output or structured JSON for programmatic integration.

## Overview

**System Health Checker** is designed to be simple, extensible, and production-ready. It performs periodic or on-demand health assessments and reports results with configurable severity thresholds. The layered architecture separates data models, metric collection, and output formatting, making it easy to test, maintain, and extend.

### Key Characteristics

- **Cross-platform**: Uses `gopsutil` for compatibility across Linux, macOS, and Windows
- **Zero dependencies for core logic**: Output formatting and metric collection are independent
- **Exit codes for automation**: Returns meaningful codes (`0`, `1`, `2`, `3`) for integration with monitoring systems and shell scripts
- **Configurable thresholds**: Override defaults via CLI flags for warning and critical levels
- **Two output formats**: Pretty-printed tables or machine-readable JSON
- **Optional process monitoring**: Check CPU and memory usage for specific processes by name

## Features

### Metrics Collected

- **CPU Usage**: Total system CPU utilization (percentage)
- **Memory Usage**: Used and total memory with percentage calculation
- **Disk Usage**: Per-mount-point disk consumption (used/total bytes and percentages)
- **Process Monitoring** (optional): PID, memory percentage, and status for a named process

### Output Formats

1. **Table (default)**: Color-coded terminal output with status indicators (🟢 OK, 🟡 WARNING, 🔴 CRITICAL)
2. **JSON**: Timestamped, structured output for scripting and tool integration

### Health Status Determination

Each metric is evaluated against configurable thresholds:
- **CRITICAL**: Metric value exceeds the critical threshold
- **WARNING**: Metric value exceeds the warning threshold but is below critical
- **OK**: Metric is within acceptable range

Overall system status is the highest severity level detected.

## Project Structure

```
healthchecker/
├── main.go                          # CLI entry point and program flow
├── go.mod                           # Module definition
├── go.sum                           # Dependency checksums
├── README.md                        # This file
│
├── internal/
│   ├── models/
│   │   ├── metrics.go               # SystemMetrics and helper methods
│   │   ├── disk.go                  # DiskInfo with status and percentage helpers
│   │   ├── process.go               # ProcessInfo data structure
│   │   └── threshold.go             # Thresholds configuration and defaults
│   │
│   ├── checker/
│   │   ├── checker.go               # HealthChecker orchestrator and status determination
│   │   ├── cpu.go                   # CPU usage collection via gopsutil
│   │   ├── memory.go                # Memory usage collection via gopsutil
│   │   ├── disk.go                  # Disk usage collection per partition
│   │   └── process.go               # Process lookup and metrics collection
│   │
│   └── output/
│       ├── table.go                 # Terminal table rendering with color
│       └── json.go                  # JSON serialization and output
│
└── test/                            # Unit tests (example test cases)
    └── checker_test.go
```

### Architecture Layers

**Models** (`internal/models/`)
- Define data structures: `SystemMetrics`, `DiskInfo`, `ProcessInfo`, `Thresholds`
- Include helper methods for calculations (e.g., percentage getters)
- No external dependencies; pure data containers

**Checker** (`internal/checker/`)
- Collects metrics using `gopsutil` library
- `HealthChecker` type orchestrates all checks
- `GetOverallStatus()` determines system health level
- Separated check functions for maintainability

**Output** (`internal/output/`)
- Format and render collected metrics
- `PrintTable()` for human-readable terminal output with colors
- `PrintJSON()` for structured data export
- Can be extended with additional formats (Prometheus, InfluxDB, etc.)

## Installation

### Prerequisites

- Go 1.23 or later (see `go.mod`)
- Support for Linux, macOS, or Windows

### Build

```bash
# Clone the repository
git clone https://github.com/yourusername/system-health-checker.git
cd healthchecker

# Build the executable
go build -o healthchecker .

# Or install directly to your $GOPATH/bin
go install .
```

### Verify Installation

```bash
./healthchecker --help
```

## Usage

### Basic Commands

```bash
# Display system health as colored table (default)
./healthchecker

# Output as JSON
./healthchecker -format=json

# Monitor a specific process
./healthchecker -process=nginx

# Use custom thresholds
./healthchecker -cpu-warning=70 -cpu-critical=85 -mem-warning=70
```

### CLI Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-format` | string | `table` | Output format: `table` or `json` |
| `-process` | string | `` | (Optional) Monitor specific process by name |
| `-cpu-warning` | float64 | `80.0` | CPU warning threshold (percent) |
| `-cpu-critical` | float64 | `90.0` | CPU critical threshold (percent) |
| `-mem-warning` | float64 | `75.0` | Memory warning threshold (percent) |
| `-mem-critical` | float64 | `85.0` | Memory critical threshold (percent) |
| `-disk-warning` | float64 | `20.0` | Disk warning threshold (percent free) |
| `-disk-critical` | float64 | `10.0` | Disk critical threshold (percent free) |

All thresholds are optional; omit the flag to use the default.

### Exit Codes

| Code | Meaning | Use Case |
|------|---------|----------|
| `0` | OK | All metrics within acceptable range |
| `1` | WARNING | At least one metric exceeded warning threshold |
| `2` | CRITICAL | At least one metric exceeded critical threshold |
| `3` | ERROR | Initialization, configuration, or runtime failure |

### Examples

#### Example 1: Basic Health Check
```bash
$ ./healthchecker
╔════════════════════════════════════════════════════════════════╗
║   SYSTEM HEALTH CHECK REPORT                                   ║
║   2025-01-15 10:30:45                                          ║
╚════════════════════════════════════════════════════════════════╝
Metric          Value                              Status     Threshold
------          -----                              ------     ---------
CPU Usage       45.20%                             ✅ OK      < 80%
Memory Usage    5.25GB / 16.00GB (32.8%)           ✅ OK      < 75%
Disk /          250.50GB / 500.00GB (50.1% used)   ✅ OK      < 20% free
Disk /home      180.75GB / 1000.00GB (18.1% used) ⚠️  WARNING < 20% free

Overall Status: ⚠️  WARNING
```

#### Example 2: JSON Output
```bash
$ ./healthchecker -format=json
{
  "timestamp": "2025-01-15T10:30:45Z",
  "overall_status": "OK",
  "metrics": {
    "cpu": {
      "percent": 45.2,
      "status": "OK"
    },
    "memory": {
      "used_bytes": 5637144576,
      "total_bytes": 17179869184,
      "percent": 32.8,
      "status": "OK"
    },
    "disks": [
      {
        "mount_point": "/",
        "used_bytes": 268435456000,
        "total_bytes": 536870912000,
        "used_percent": 50.0,
        "free_percent": 50.0,
        "status": "OK"
      },
      {
        "mount_point": "/home",
        "used_bytes": 194340020224,
        "total_bytes": 1099511627776,
        "used_percent": 18.0,
        "free_percent": 82.0,
        "status": "OK"
      }
    ],
    "processes": []
  }
}
```

#### Example 3: Process Monitoring
```bash
$ ./healthchecker -process=nginx -format=json
{
  "timestamp": "2025-01-15T10:32:10Z",
  "overall_status": "OK",
  "metrics": {
    ...
    "processes": [
      {
        "name": "nginx",
        "pid": 1234,
        "status": "running",
        "memory_percent": 2.5
      }
    ]
  }
}
```

#### Example 4: Custom Thresholds
```bash
$ ./healthchecker \
  -cpu-warning=70 \
  -cpu-critical=85 \
  -mem-warning=70 \
  -mem-critical=80 \
  -disk-warning=15 \
  -disk-critical=5
```

## Thresholds and Status Calculation

### Default Thresholds

| Metric | Warning | Critical | Notes |
|--------|---------|----------|-------|
| CPU | 80% | 90% | Percentage of total CPU used |
| Memory | 75% | 85% | Percentage of total memory used |
| Disk | 20% free | 10% free | Percentage of free space remaining |

### Status Determination Logic

**For CPU and Memory** (higher values indicate worse health):
```
if value >= critical_threshold:
  status = CRITICAL
else if value >= warning_threshold:
  status = WARNING
else:
  status = OK
```

**For Disk** (uses free percent; lower free space is worse):
```
free_percent = 100 - used_percent
if free_percent < critical_threshold:
  status = CRITICAL
else if free_percent < warning_threshold:
  status = WARNING
else:
  status = OK
```

**Overall Status** (worst severity wins):
```
if any_metric == CRITICAL:
  overall = CRITICAL
else if any_metric == WARNING:
  overall = WARNING
else:
  overall = OK
```

## Output Formats

### Table Format

**Strengths:**
- Human-readable and color-coded for quick visual assessment
- Shows threshold information for context
- Emoji indicators for instant status recognition

**Details:**
- Uses `github.com/fatih/color` for terminal colors
- Uses `github.com/olekukonko/tablewriter` for structured table rendering
- Status colors: 🟢 (green) for OK, 🟡 (yellow) for WARNING, 🔴 (red) for CRITICAL

### JSON Format

**Strengths:**
- Structured and machine-parseable
- Timestamp included for correlation with other events
- Suitable for logs, metrics collection, and API integration
- Omits empty process array to reduce payload size

**Structure:**
```json
{
  "timestamp": "ISO 8601 timestamp",
  "overall_status": "OK|WARNING|CRITICAL",
  "metrics": {
    "cpu": {
      "percent": number,
      "status": "OK|WARNING|CRITICAL"
    },
    "memory": {
      "used_bytes": integer,
      "total_bytes": integer,
      "percent": number,
      "status": "OK|WARNING|CRITICAL"
    },
    "disks": [
      {
        "mount_point": "string",
        "used_bytes": integer,
        "total_bytes": integer,
        "used_percent": number,
        "free_percent": number,
        "status": "OK|WARNING|CRITICAL"
      }
    ],
    "processes": [
      {
        "name": "string",
        "pid": integer,
        "status": "string",
        "memory_percent": number
      }
    ]
  }
}
```

## Testing

### Run Tests

```bash
go test ./...          # Run all tests
go test -v ./...       # Verbose output
go test -race ./...    # Detect race conditions
```

### Test Coverage

Generate and view test coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Example Unit Tests

The `test/` directory includes examples:
- `DiskInfo.GetUsedPercent()` calculation verification
- `DiskInfo.GetStatus()` threshold evaluation
- `SystemMetrics.GetMemoryPercent()` computation

## Extending the Project

### Adding a New Check

1. **Create a new checker file** in `internal/checker/`:
   ```go
   // internal/checker/network.go
   package checker
   
   func (hc *HealthChecker) CheckNetwork() error {
       // Collect network metrics using gopsutil
       // Populate hc.metrics with results
       return nil
   }
   ```

2. **Update `CheckAll()`** in `internal/checker/checker.go`:
   ```go
   func (hc *HealthChecker) CheckAll() error {
       // ... existing checks ...
       if err := hc.CheckNetwork(); err != nil {
           return fmt.Errorf("network check failed: %w", err)
       }
       return nil
   }
   ```

3. **Update model** in `internal/models/metrics.go` if new fields are needed

4. **Update output formatters** in `internal/output/`:
   - Add fields to JSON structs
   - Add rows to table output

### Adding a New Output Format

1. **Create formatter file** in `internal/output/`:
   ```go
   // internal/output/prometheus.go
   package output
   
   func PrintPrometheus(metrics *models.SystemMetrics, thresholds *models.Thresholds) {
       // Format metrics as Prometheus exposition format
   }
   ```

2. **Wire into `main.go`**:
   ```go
   case "prometheus":
       output.PrintPrometheus(metrics, thresholds)
   ```

### Integration Examples

**Send metrics to a monitoring system:**
```bash
./healthchecker -format=json | curl -X POST \
  -H "Content-Type: application/json" \
  -d @- https://monitoring.example.com/api/metrics
```

**Log output with timestamp:**
```bash
./healthchecker >> /var/log/healthchecker.log 2>&1
```

**Conditional alerting:**
```bash
if ./healthchecker > /dev/null 2>&1; then
  status=$?
  case $status in
    2) echo "CRITICAL" | mail -s "System Critical" ops@example.com ;;
    1) echo "WARNING" | logger -p local0.warning ;;
  esac
fi
```

## Implementation Details

### Metric Collection

- **CPU**: Uses `github.com/shirou/gopsutil/v4/cpu.Percent()` with interval=0 for instantaneous reading
- **Memory**: Uses `github.com/shirou/gopsutil/v4/mem.VirtualMemory()` for system memory stats
- **Disk**: Uses `github.com/shirou/gopsutil/v4/disk.Partitions()` and `disk.Usage()` per mount point
- **Process**: Uses `github.com/shirou/gopsutil/v4/process.Processes()` and name matching

### Conversions

- **Bytes to GB**: Divides by 1,024³ (1024 × 1024 × 1024)
- **Percentages**: Calculated as (used / total) × 100

### Color Support

Terminal color output requires ANSI-compatible terminal. On Windows, uses `github.com/mattn/go-colorable` for compatibility.

## Performance Considerations

- **CPU sampling**: `interval=0` provides instant snapshot; consider non-zero interval for smoother readings
- **Memory**: Single system call, minimal overhead
- **Disk**: Iterates all mounted partitions; may vary based on system configuration
- **Process lookup**: Scans all processes; slower on systems with many processes

## Troubleshooting

### Common Issues

**Problem**: Process not found
```
process check warning: process not found: nginx
```
**Solution**: Verify the process is running and use the exact executable name.

**Problem**: Permission denied
```
health checks failed: open /sys/...: permission denied
```
**Solution**: Some metrics require elevated privileges. Run with `sudo` if necessary.

**Problem**: JSON output with special characters
**Solution**: JSON is properly escaped; use `json` tools for parsing (e.g., `jq`).

## Development

### Building from Source

```bash
git clone https://github.com/yourusername/system-health-checker.git
cd healthchecker
go mod download
go build -o healthchecker .
```

### Code Style

- Follow standard Go conventions
- Use `go fmt`, `go vet`, and `golint`
- Write comments for exported functions

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature/your-feature`
5. Submit a pull request with description and test coverage

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/shirou/gopsutil/v4` | v4.25.11 | System metrics collection |
| `github.com/fatih/color` | v1.18.0 | Terminal color output |
| `github.com/olekukonko/tablewriter` | v1.1.2 | Table formatting |

All dependencies are managed by Go modules. Run `go mod tidy` to clean up.

## License

This project is released as-is. Add your preferred license file (`LICENSE`) to the repository root based on your organization's policies.

## Support and Contact

For issues, questions, or contributions:
- Open an issue on GitHub
- Submit a pull request
- Review `CONTRIBUTING.md` for guidelines (if available)

---

**Happy monitoring!** 🚀