package abi

// Example of real JSON: see below
type paramType string

type jsonParam struct {
	Name string    `json:"name"`
	Type paramType `json:"type"`
}

const (
	tFunction paramType = "function"
	//tUint256  paramType = "uint256"
	//tAddress  paramType = "address"
)

type jsonItem struct {
	Name            string      `json:"name"`
	Constant        bool        `json:"constant"`
	Payable         bool        `json:"payable"`
	StateMutability string      `json:"stateMutability"`
	Type            paramType   `json:"type"`
	Inputs          []jsonParam `json:"inputs"`
	Outputs         []jsonParam `json:"outputs"`
}

/*
Example of real JSON:
[
   {
      "constant":false,
      "inputs":[
         {
            "name":"_spender",
            "type":"address"
         },
         {
            "name":"_value",
            "type":"uint256"
         }
      ],
      "name":"approve",
      "outputs":[
         {
            "name":"success",
            "type":"bool"
         }
      ],
      "payable":false,
      "stateMutability":"nonpayable",
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[

      ],
      "name":"transfersEnabled",
      "outputs":[
         {
            "name":"",
            "type":"bool"
         }
      ],
      "payable":false,
      "stateMutability":"view",
      "type":"function"
   },
   {
      "constant":true,
      "inputs":[
         {
            "name":"",
            "type":"address"
         },
         {
            "name":"",
            "type":"uint256"
         }
      ],
      "name":"agingBalanceOf",
      "outputs":[
         {
            "name":"",
            "type":"uint256"
         }
      ],
      "payable":false,
      "stateMutability":"view",
      "type":"function"
   },
	{
      "anonymous":false,
      "inputs":[
         {
            "indexed":true,
            "name":"from",
            "type":"address"
         },
         {
            "indexed":false,
            "name":"value",
            "type":"uint256"
         }
      ],
      "name":"Burn",
      "type":"event"
   },
   ...
]
*/
