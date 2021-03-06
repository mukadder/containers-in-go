package funding

type FundServer struct {
	Commands chan interface{}
	fund     Fund
}

func NewFundServer(initialBalance int) *FundServer {
	server := &FundServer{
		// make() creates builtins like channels, maps, and slices
		Commands: make(chan interface{}),
		fund:     NewFund(initialBalance),
	}

	// Spawn off the server's main loop immediately
	go server.loop()
	return server
}

ffunc (s *FundServer) loop() {
    for command := range s.Commands {

        // command is just an interface{}, but we can check its real type
        switch command.(type) {

        case WithdrawCommand:
            // And then use a "type assertion" to convert it
            withdrawal := command.(WithdrawCommand)
            s.fund.Withdraw(withdrawal.Amount)

        case BalanceCommand:
            getBalance := command.(BalanceCommand)
            balance := s.fund.Balance()
            getBalance.Response <- balance

        default:
            panic(fmt.Sprintf("Unrecognized command: %v", command))
        }
    }
}
}
