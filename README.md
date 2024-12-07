# Blackjack-Lab

_A Golang-based platform for simulating blackjack strategies, evaluating Expected Value (EV), and managing portfolios._

---

## Features

- **Blackjack Simulation**:
  - Supports 6-deck shuffling for realistic gameplay.
  - Customize strategies for evaluating Expected Value (EV).
- **Portfolio Management**:
  - Simulate multiple strategies across a timeframe for comparative analysis.
  - Track profits, losses, and EV trends.
- **Advanced Mechanics**:
  - Implements blackjack core components: hits, stands, doubles, splits, etc.
  - Detailed game logic with robust probability models.
- **Efficient Processing**:
  - Optimized for large-scale simulations, handling thousands of iterations.
- **Extendable**:
  - Future support planned for additional variants of blackjack and extended card counting systems.
- **Golang-based**:
  - High performance with clean, modular architecture.

---

## Getting Started

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jamie0xgitc0decat/blackjack-lab.git
   ```
2. **Navigate to the project directory:**
   ```bash
   cd blackjack-lab
   ```
3. **Build the project:**
   ```bash
   go build
   ```
4. **Run the application:**
   ```bash
   ./blackjack-lab
   ```

---

## Documentation

### Core Components

- **Deck Management**:
  - Simulates a 6-deck shoe with realistic shuffle algorithms.
- **Strategy Customization**:
  - Define and apply strategies such as:
    - Basic Strategy
    - Custom EV Maximization Strategies
- **Simulation Engine**:
  - Run simulations over thousands of hands or even entire months.
  - Analyze performance metrics for profitability.

### Supported Strategies

- **Basic Strategy**:
  - Optimal decision-making for standard rules.
- **Custom Strategies**:
  - Define unique rules tailored for specific card-counting methods or risk profiles.

---

## Features Roadmap

- [x] Support for 6-deck blackjack simulations
- [x] Custom strategy configuration

---

## Contributing
welcome contributions! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Acknowledgments

- Inspired by **Blackjack enthusiasts**.


---

