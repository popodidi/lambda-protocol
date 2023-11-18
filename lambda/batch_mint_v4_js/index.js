const main = (obj) => {
  const { addr } = obj
  // get the lest 40 chars of the address
  const addrEnd = addr.slice(-40)

  return [
    {
      "data":"0x",
      "to":`0x${addrEnd}`,
      "value":"1000000000",
      "from":"0x5bfd83c76641bea157d670e8e0d42a46ac144298",
      "txHash":"0xbe19ab9e5031bb20e8e4b7e2da795ebccb594506c284fde4752d0e240c3d5519",
      "methodData":{
        "name":"transfer",
        "params":[]
      },
      "readableCallData":"transfer"
    },
    {
      "data":"0xa0712d6800000000000000000000000000000000000000000000000000000000000000010021fb3f",
      "to":"0xc927ebf64654d3aebc967be54b93a330decdcd58",
      "value":"0",
      "from":"0x5bfd83c76641bea157d670e8e0d42a46ac144298",
      "txHash":"0xee937d1be0f324a86fd873f4ad1c0e1c277c4fa2eb1f180e7f7a5ecfb16d6975",
      "methodData":{"name":"mint","params":[{"name":"quantity","value":"1","type":"uint256"}]},
      "readableCallData":"mint"
    },
    {
      "data":"0xef4c6687000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001568747470733a2f2f65746865727363616e2e696f2f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021fb3f",
      "to":"0x48fbf0b7e6bdaa9d7452feb34507b2bddb148e54",
      "value":"0",
      "from":"0x5bfd83c76641bea157d670e8e0d42a46ac144298",
      "txHash":"0x920ae9fe07d69b17dfdbc9c718b96b61e0769c869c40f61ae0d4ee873513536b",
      "methodData":{
        "name":"mintSlug",
        "params":[{"name":"url","value":"https://etherscan.io/","type":"string"},{"name":"slug","value":"","type":"string"},{"name":"referrer","value":"0x0000000000000000000000000000000000000000","type":"address"}]
      },
      "readableCallData":"mintSlug"
    }
  ]
}
