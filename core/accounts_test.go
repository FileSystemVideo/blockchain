package core

import "testing"

func TestAa1(t *testing.T) {
	/**
	: fsv17xpfvakm2amg962yls6f84z3kell8c5lmrfnas
	:  fsv1gwqac243g2z3vryqsev6acq965f9ttwhjpnmg9
	:  fsv1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8v9w0lj
	: fsv1fl48vsnmsdzcv85q5d2q4z5ajdha8yu37pru06
	: fsv1tygms3xhhs3yv487phx3dw4a95jn7t7l2pldew
	*/
	
	t.Log(":", ContractAddressBonus.String())

	t.Log(":", ContractAddressFee.String())

	t.Log(":", ContractAddressBank.String()) //fsv1gwqac243g2z3vryqsev6acq965f9ttwhjpnmg9

	t.Log(":", ContractAddressDistribution.String())

	t.Log(":", ContractAddressDeflation.String())

	t.Log(":", ContractAddressStakingBonded.String())

	t.Log(":", ContractAddressStakingNotBonded.String())
	t.Log("", ContractAddressGov.String())
}
