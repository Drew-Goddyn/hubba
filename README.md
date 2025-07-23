# Terminal Application with Entity Management and Animation

A terminal-based application built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) that implements entity management, collision detection, and smooth animations.

![Go Version](https://img.shields.io/badge/Go-1.23.4+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Cross--Platform-orange.svg)

## Technical Features

### Entity Management System
- **Dynamic entity creation** with configurable properties
- **Multiple entity types**: Spheres and sprites with distinct behaviors
- **Animation system** using Harmonica library for smooth movement
- **Parameter adjustment** for gravity, bounce coefficients, size, and colors

### User Interface
- **Split-screen layout** with responsive design
- **Adaptive interface** scaling from 50+ character terminals to ultra-wide displays
- **Styled components** using Lip Gloss for visual rendering
- **Interactive controls** with keyboard and button navigation
- **Real-time metrics display** including FPS and entity count

### Technical Architecture
- **Polymorphic entity system** with interface-based design
- **Modular calculation engine** with configurable parameters
- **Comprehensive test suite** with >95% code coverage
- **Performance optimization** for 60 FPS rendering
- **Stress testing** with configurable entity limits

## Interface Examples

### Standard Operation Mode
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│  ●   ○   ■   ●   ■      ●○●   ●   ○   ■      ●●   ●   ◆   ●   ○   ●   ○         │
│ ○       ◆       ●◆     ●○    ■   ●     ■        ●   ◆   ○   ●●                  │
│                                                                                 │
│                                                                                 │
│ ⚙️ Gravity: 25.0 | 🏀 Bounce: 0.70 | 📊 FPS: 60.0 | 🎯 Limit: 1000          │
│                                                                                 │
│ Entities: 40 FPS: 60.0                                                         │
└─────────────────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           🎮 PHYSICS CONTROLS                                  │
│ →Add Sphere← Add Sprite Clear All Pause Reset                                  │
│ ⚙️Normal 📏Tiny 🎨Spring Green                                                  │
│ Keys: A=Add●  S=Add◆  C=Clear  P=Pause  R=Reset  G=Gravity  B=Bounce  Z=Size  │
│ X=Color  F=Perf  T=Test  L=Limit  TAB=Navigate                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Performance Monitoring Mode
```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                                                                 │
│  ●   ○   ■   ●   ■      ●○●   ●   ○   ■      ●●   ●   ◆   ●   ○   ●   ○         │
│ ○       ◆       ●◆     ●○    ■   ●     ■        ●   ◆   ○   ●●                  │
│                                                                                 │
│ ⚙️ Gravity: 25.0 | 🏀 Bounce: 0.70 | 📊 FPS: 60.0 | 🎯 Limit: 1000          │
│ 📐 Terminal: 103x46 | Sim: 97x32 | Ctrl: 97x6                                 │
│                                                                                 │
│ Entities: 40 FPS: 60.0                                                         │
└─────────────────────────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           🎮 PHYSICS CONTROLS                                  │
│ →Add Sphere← Add Sprite Clear All Pause Reset                                  │
│ ⚙️Normal 📏Tiny 🎨Spring Green                                                  │
│ Keys: A=Add●  S=Add◆  C=Clear  P=Pause  R=Reset  G=Gravity  B=Bounce  Z=Size  │
│ X=Color  F=Perf  T=Test  L=Limit  TAB=Navigate                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### Compact Terminal Mode (50-80 characters)
```
┌──────────────────────────────────────────────────┐
│                                                  │
│ ●  ○  ■     ●○●  ●    ○  ■     ●●  ●  ◆  ●  ○    │
│○     ◆      ●◆   ●○   ■  ●      ■    ●  ◆  ○     │
│                                                  │
│ ⚙️ Gravity: 25.0 | 🏀 Bounce: 0.70              │
│ Entities: 40 FPS: 60.0                          │
└──────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────┐
│                PHYSICS CONTROLS                  │
│ →Add Sphere← Add Sprite Clear All Pause Reset    │
│ ⚙️Normal 📏Tiny 🎨Spring                         │
│ Keys: A●  S◆  C=Clear  P=Pause  F=Perf          │
└──────────────────────────────────────────────────┘
```

### Key Interface Elements

- **Upper Panel**: Entity display area with real-time positioning
- **Lower Panel**: Interactive controls and parameter settings
- **Entity Symbols**: Circles (●○), squares (■), diamonds (◆) in various colors
- **Status Bar**: Physics parameters, FPS counter, entity count
- **Control Buttons**: Keyboard shortcuts and parameter displays
- **Responsive Layout**: Adapts to terminal width (50+ to 200+ characters)

## Installation and Setup

### Prerequisites
- **Go 1.23.4+** installed
- Terminal with minimum **50 characters width**

### Installation Steps

1. **Clone repository:**
   ```bash
   git clone <repository-url>
   cd bubblegum-test
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Build application:**
   ```bash
   go build -o bubblegum-physics-sim .
   ```

4. **Execute:**
   ```bash
   ./bubblegum-physics-sim
   ```

   Or run directly:
   ```bash
   go run .
   ```

## Operation Controls

### Entity Management
| Key | Function | Description |
|-----|----------|-------------|
| `a` | Add Sphere | Creates sphere entity with current parameters |
| `s` | Add Sprite | Creates sprite entity with current parameters |
| `c` | Clear All | Removes all entities |
| `p` | Pause/Resume | Toggles entity updates (animations continue) |
| `r` | Reset | Clears entities and resumes updates |
| `q` / `Ctrl+C` | Exit | Terminates application |

### Parameter Controls
| Key | Parameter | Values |
|-----|-----------|--------|
| `g` | Gravity | Zero → Low → Normal → High |
| `b` | Bounce | No Bounce → Low → Normal → Perfect |
| `z` | Entity Size | Tiny → Small → Medium → Large |
| `x` | Entity Color | Cycles through 16 colors |

### System Controls
| Key | Feature | Description |
|-----|---------|-------------|
| `f` | Performance Mode | Shows FPS, limits, and system info |
| `t` | Stress Test | Adds 20 entities rapidly |
| `l` | Entity Limit | Cycles between 1000 → 2000 → 5000 limits |

### Navigation
- **Tab/Shift+Tab**: Navigate UI elements
- **Enter/Space**: Activate focused element
- **Arrow Keys**: Navigate controls
- **Mouse**: Direct element interaction (where supported)

## Code Architecture

### Core Components

#### Entity System (`entities.go`)
```go
type Entity interface {
    GetPosition() (float64, float64)
    SetPosition(x, y float64)
    GetVelocity() (float64, float64)
    SetVelocity(vx, vy float64)
    // ... additional methods
}
```

**Implementations:**
- **Spheres**: Circular entities with size-based properties
- **Sprites**: Customizable entities with symbol representation

#### Calculation Engine (`physics.go`)
```go
type PhysicsEngine struct {
    Gravity       float64 // Gravity acceleration
    AirResistance float64 // Air resistance coefficient
    Restitution   float64 // Bounce factor
    // ... additional properties
}
```

**Functions:**
- Gravity application
- Collision detection and response
- Velocity calculations
- Boundary enforcement

#### Animation System (`animation.go`)
```go
type AnimationEngine struct {
    SpringTension float64 // Spring stiffness
    SpringDamping float64 // Spring damping
    TargetFPS     int     // Target frame rate
}
```

**Capabilities:**
- Spring-based interpolation
- 60 FPS rendering target
- Position smoothing
- Harmonica integration

#### Control Interface (`controls.go`)
```go
type ControlPanel struct {
    buttons      []Button
    focused      int
    buttonStyles ButtonStyles
    // ... layout properties
}
```

**Features:**
- Responsive layouts
- Input handling
- Parameter display
- Adaptive sizing

### File Structure

```
bubblegum-test/
├── main.go              # Application entry and UI
├── entities.go          # Entity interfaces and types
├── physics.go           # Calculation engine
├── animation.go         # Animation system
├── controls.go          # User interface controls
│
├── *_test.go           # Test suite
├── TEST_COVERAGE.md    # Coverage documentation
│
├── go.mod              # Module definition
├── go.sum              # Dependency hashes
└── README.md           # Documentation
```

## Testing

### Test Execution
```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Categories
- **Unit Tests**: Component-level testing
- **Integration Tests**: Cross-component functionality
- **Performance Tests**: Load and stress testing
- **Edge Case Tests**: Boundary condition validation
- **UI Tests**: Interface behavior validation

### Coverage Metrics
- **Overall Coverage**: >95%
- **Critical Paths**: 100% covered
- **Details**: See `TEST_COVERAGE.md`

## Performance Characteristics

### Optimization Features
- **60 FPS rendering** with adaptive timing
- **Spatial optimization** for collision detection
- **Conditional evaluation** for off-screen entities
- **Memory management** for frequent allocations
- **Responsive scaling** based on terminal dimensions

### Performance Data
| Load Level | Entity Count | FPS Range | Memory Usage |
|------------|--------------|-----------|--------------|
| Light | 1-20 | 60+ | <10MB |
| Medium | 21-100 | 45-60 | 10-25MB |
| Heavy | 101-500 | 30-45 | 25-50MB |
| Stress | 500-1000 | 15-30 | 50-100MB |

### Performance Monitoring
1. **Enable Performance Mode** (`f` key) for real-time metrics
2. **Adjust entity limits** (`l` key) based on system capacity
3. **Smaller terminal windows** reduce computational load
4. **Pause functionality** (`p` key) reduces CPU usage

## Configuration

### Color Options
16 available colors:
- Spring Green, Gold, Dodger Blue, Hot Pink
- Dark Turquoise, Orange Red, Alice Blue, Tomato
- Turquoise, Sky Blue, Pale Green, Orange
- Orchid, Light Sea Green, Light Pink, Green Yellow

### Engine Parameters
```go
// Configurable values in physics.go
Gravity:       25.0  // Acceleration (units/second²)
AirResistance: 0.05  // Drag coefficient (0-1)
Restitution:   0.7   // Bounce factor (0-1)
StaticFriction: 0.8  // Friction coefficient (0-1)
```

### UI Breakpoints
Terminal width handling in `main.go`:
- **Ultra-compact**: <50 characters (minimal interface)
- **Compact**: 50-80 characters (essential features)
- **Standard**: 80-120 characters (full features)
- **Wide**: >120 characters (enhanced visuals)

## Build Configuration

### Build Commands
```bash
# Development build
go build -o bubblegum-physics-sim .

# Optimized build
go build -ldflags="-s -w" -o bubblegum-physics-sim .

# Cross-compilation (Linux example)
GOOS=linux GOARCH=amd64 go build -o bubblegum-physics-sim-linux .
```

### Development Process
1. **Execute tests**: `go test -v`
2. **Check coverage**: `go test -cover`
3. **Format code**: `go fmt ./...`
4. **Lint code**: `golangci-lint run` (if available)
5. **Build and test**: `go build . && ./bubblegum-physics-sim`

### Adding Entity Types
```go
// 1. Implement Entity interface
type CustomEntity struct {
    BaseEntity
    property string
}

// 2. Add to EntityType constants
const CustomEntityType EntityType = "custom"

// 3. Add factory method
func NewCustomEntity(x, y float64) *CustomEntity {
    // Implementation
}
```

## Troubleshooting

### Build Issues
**Go version errors**
```
Solution: Install Go 1.23.4+
- Update: brew upgrade go (macOS)
- Verify: go version
- Clear: go clean -cache
```

### Performance Issues
**Low FPS or stuttering**
```
Solutions:
1. Reduce entity count ('l' key)
2. Use smaller terminal
3. Enable performance mode ('f' key)
4. Close other applications
```

### Display Issues
**Layout problems or text overlap**
```
Solutions:
1. Ensure terminal width ≥50 characters
2. Try different terminal emulator
3. Check color support
4. Restart application
```

### Input Issues
**Unresponsive keys**
```
Solutions:
1. Ensure terminal focus
2. Try Alt+Screen toggle
3. Check key conflicts
4. Restart with 'q' and relaunch
```

### System Monitoring
Performance mode (`f` key) displays:
- Real-time FPS
- Entity count and limits
- Memory indicators
- Terminal dimensions
- Parameter values

## API Reference

### Entity Interface
```go
type Entity interface {
    // Position management
    GetPosition() (float64, float64)
    GetDisplayPosition() (float64, float64)
    SetPosition(x, y float64)
    SetImmediatePosition(x, y float64)
    
    // Dynamics
    GetVelocity() (float64, float64)
    SetVelocity(vx, vy float64)
    ApplyForce(fx, fy float64)
    Update(deltaTime float64)
    
    // Visual properties
    GetSymbol() string
    GetColor() lipgloss.Color
    GetSize() int
    Render() string
    
    // Collision and animation
    GetBounds() (x, y, width, height float64)
    CheckCollision(other Entity) bool
    UpdateAnimation(ae *AnimationEngine)
}
```

### Engine Methods
```go
func (pe *PhysicsEngine) ApplyPhysics(entities []Entity)
func (pe *PhysicsEngine) HandleEntityCollisions(entities []Entity)
func (pe *PhysicsEngine) SetGravity(gravity float64)
func (pe *PhysicsEngine) SetRestitution(restitution float64)
func (pe *PhysicsEngine) AddRandomVelocity(entity Entity, maxSpeed float64)
```

## Development Guidelines

### Code Standards
1. **Follow Go conventions** and use `go fmt`
2. **Write tests** for new functionality
3. **Update documentation** for API changes
4. **Test multiple terminal sizes** (50-200+ chars)
5. **Maintain >95% test coverage**

### Project Management
Uses [Task Master AI](https://github.com/codeium/taskmaster):
- See `.taskmaster/` directory for task tracking
- Automated task generation from requirements
- Progress and dependency management

## License

MIT License - see [LICENSE](LICENSE) file.

## Dependencies

- **[Charm](https://charm.sh/)** - TUI ecosystem
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling
- **[Harmonica](https://github.com/charmbracelet/harmonica)** - Animation library
- **Task Master AI** - Project management

---

### Execute Application

```bash
go run .
``` 