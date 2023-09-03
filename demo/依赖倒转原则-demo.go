package main

/**
模拟组装两台计算机。
	抽象层:
		有显卡Card方法display，有内存Memory方法storage，有处理器CPU方法calculate。
	实现层:
		Intel(英特尔)公司，产品有显卡、内存、CPU;
		Kingston公司，产品有内存;
		NVIDIA公 司，产品有显卡。
	逻辑层:
		(1)组装一台Intel系列的计算机，并运行。
		(2)组装一台Intel CPU、Kingston内存、NVIDIA显卡的计算机，并运行。
**/

type Card interface {
	display()
}

type Memory interface {
	storage()
}

type CPU interface {
	calculate()
}

type IntelCard struct {
}

func (intelCard *IntelCard) display() {
	println("IntelCard display")
}

type IntelMemory struct {
}

func (intelMemory *IntelMemory) storage() {
	println("IntelMemory storage")
}

type IntelCPU struct {
}

func (intelCPU *IntelCPU) calculate() {
	println("IntelCPU calculate")
}

type KingstonMemory struct {
}

func (kingstonMemory *KingstonMemory) storage() {
	println("KingstonMemory storage")
}

type NVIDIA struct {
}

func (nvidia *NVIDIA) display() {
	println("NVIDIA display")
}

type Computer struct {
	card   Card
	memory Memory
	cpu    CPU
}

func (computer *Computer) setCard(card Card) {
	computer.card = card
}

func (computer *Computer) setMemory(memory Memory) {
	computer.memory = memory
}

func (computer *Computer) setCPU(cpu CPU) {
	computer.cpu = cpu
}

func (computer *Computer) run() {
	computer.card.display()
	computer.memory.storage()
	computer.cpu.calculate()
}

func main() {
	computer := &Computer{}
	computer.setCard(&IntelCard{})
	computer.setMemory(&IntelMemory{})
	computer.setCPU(&IntelCPU{})
	computer.run()
	println("======================================")
	computer.setCard(&NVIDIA{})
	computer.setMemory(&KingstonMemory{})
	computer.setCPU(&IntelCPU{})
	computer.run()
}
