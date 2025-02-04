## Example Workflow

Below is an example scenario:

### TON FASTNET

LiteClient: Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5
TxChecker: Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK


1. Deploy the system to fastnet using key block 27788014 from testnet.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml deploy all -s 27788014 -w -1
```

 - lite_client deploy [link](http://109.236.91.95:8080/transaction?account=Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5&lt=868253000003&hash=B1A9BC9C8AFD6BE8A06265EEBEA2D896CCA9F4AE65862C7B9DA687B9AAFCCAC5)
 - tx_checker deploy [link](http://109.236.91.95:8080/transaction?account=Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK&lt=868258000003&hash=AF74B3B02CC644B0FE139ECAB4216A25D86231081E22371F1A0ECB406981E099)

2. send checkBlock for a block in the same epoch → OK.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5 check-block -s 27788015
```

 [link](http://109.236.91.95:8080/transaction?account=Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5&lt=868521000003&hash=71A9DDF28D25F0FF45B2E281F4E544FCC9EE405D24DA4AB63010A694ADA153C2)

3. send checkBlock for the previous epoch → Not OK.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5 check-block -s 27788013
```

[link](http://109.236.91.95:8080/transaction?account=Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5&lt=868733000003&hash=558710D4D15DF63F5ACAB9E36845F3452D81FC4D91F5F475EC7FBCB07D2CC7CD)

4. send checkTransaction with a valid transaction in the current epoch → OK.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK check-tx -s 27788020 -t 7E55E639D4EF717A89A06C5F20937768151E31F822DCFA04751C743721CCCA21
```

Contract chain: wallet -> tx_checker -> lite_client -> tx_checker -> wallet (transaction_checked#756adff1)

Link to tx from lite_client to checker:
[link](http://109.236.91.95:8080/transaction?account=Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK&lt=869561000007&hash=763A42F880BC1DB2054E8EB93A381EFED20EE34C20C5B2F16C3644B238C27E84)

5. send checkTransaction with a transaction in the prev epoch → Not OK.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK check-tx -s 27788013 -t E0F62001C2F78F3FB54199762EE382A7AF03AEF4D86A85EF6BA8FEDB58604CBF
```
Contract chain: wallet -> tx_checker -> lite_client -> (bounce) tx_checker

Link to tx with bounced msg from lite_client: 
[link](http://109.236.91.95:8080/transaction?account=Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK&lt=869370000007&hash=42EA7CEF564823FBE614288CFF7E6071D4097E6F060D35F27043FFAD377532C9)

6. send newKeyBlock with an old block → Not OK.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5 new-key-block -s 27787289
```
[link](http://109.236.91.95:8080/transaction?account=Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5&lt=869877000003&hash=0B98331D09F4724F7DAC39EAF3B8BE017288C6E3287230A9D830D0F5849D2995)

7. send newKeyBlock with new epoch -> OK.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5 new-key-block -s 27793804
```

[link](http://109.236.91.95:8080/transaction?account=Ef_cmIsszQinqjDnK4LIib3vSBE8Zhf-ytgRJDGispoD-Et5&lt=870276000003&hash=3D6C725D8F591BD63E9127D645AF9534BAF9BECF79C4C0E05481FC981DDD2001)

8. send checkTransaction with tx in new epoch -> Ok.

```bash
go run main.go --network fastnet --config ./.trustless-bridge-cli.yaml send -a Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK check-tx -s 27793805 -t E9FC1788AA8976959922CE27EDF0CB3E7622F62777EFF07037DB2099208E8BBB
```

Contract chain: wallet -> tx_checker -> lite_client -> tx_checker -> wallet

Link to tx from lite_client to tx_checker: 
[link](http://109.236.91.95:8080/transaction?account=Ef8zZWfeh22ib982EIgo_FZM0n2Iym1WHFBRBA_H7BEfsoMK&lt=870717000007&hash=45E746F7661EB58A69A3A0F99DB6AF75E95ED7B3FE442AA0F37655C4E370E3FC)

### TON TESTNET
LiteClient: EQBGWoImJJ8Uw4Lz0b2yXjpOf31awQfXHJthrYB4zppnL3c1
TxChecker: EQCEILr1N8ey9Ar-9OtnCq8A4v217lsE0pJuEZgrZOStAVVa

1. Deploy the system to testnet using key block 765944 from fastnet.

```bash
go run main.go --network testnet --config ./.trustless-bridge-cli.yaml deploy all -s 765944 -w 0
```

 - lite_client deploy [link](https://testnet.tonviewer.com/transaction/7ce957f4a9c6066809954ad204a3e80625f8193a51dc32352103e74f4cd77830)
 - tx_checker deploy [link](https://testnet.tonviewer.com/transaction/2b1c64e38044c7e53579ebf6c6c542f0bb2312243ebf8faab17dc1ddb97a8071)

2. send checkBlock for a block in the same epoch → OK.

```bash
go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQBGWoImJJ8Uw4Lz0b2yXjpOf31awQfXHJthrYB4zppnL3c1 check-block -s 765945
```

 [link](https://testnet.tonviewer.com/transaction/c615ee86016ba15f1d22a6f3690156551c037003b36b943537d3617d654d2773)

3. send checkBlock for the previous epoch → Not OK.

```bash
go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQBGWoImJJ8Uw4Lz0b2yXjpOf31awQfXHJthrYB4zppnL3c1 check-block -s 765940
```

[link](https://testnet.tonviewer.com/transaction/a510e919c39dddf13f6015a4d94972303bda0794713e468952d834d3612d3c39)

4. send checkTransaction with a valid transaction in the current epoch → OK.

```bash
go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQCEILr1N8ey9Ar-9OtnCq8A4v217lsE0pJuEZgrZOStAVVa check-tx -s 765950 -t 5052521713ECA40FF82B354FC520D026B94601660A67574A7AF0E3DA8BF68C67
```

[link](https://testnet.tonviewer.com/transaction/c017d0d845bdc0a98f895b49c1934fcab66e443a78e2d07f7f48ae3e51159852)

5. send checkTransaction with a transaction in the prev epoch → Not OK.

```bash
    go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQCEILr1N8ey9Ar-9OtnCq8A4v217lsE0pJuEZgrZOStAVVa check-tx -s 765941 -t 2FEFFECAD8215086C991F68D3C42BC57A7A801A5BB45CF720FE6C4221AA42FD0
```

[link](https://testnet.tonviewer.com/transaction/2e21908bd82f615bc38b18025a48f4f81e174e2c49cc97bb7fa66d01759635ea)

6. send newKeyBlock that does not change the epoch → OK.

```bash
    go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQBGWoImJJ8Uw4Lz0b2yXjpOf31awQfXHJthrYB4zppnL3c1 new-key-block -s 850663
```

[link](https://testnet.tonviewer.com/transaction/271c57f53c8f116f5c11d8fbb160e85e774b023ed23b4c9b80d22ae562308fbd)

7. send newKeyBlock with an old block → Not OK.

```bash
    go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQBGWoImJJ8Uw4Lz0b2yXjpOf31awQfXHJthrYB4zppnL3c1 new-key-block -s 765652
```
[link](https://testnet.tonviewer.com/transaction/2bdd6af42baed97781cd622927cf27b243083e1b2a7767e08c819075a8420242)

8. send newKeyBlock with new epoch -> OK.

```bash
go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQBGWoImJJ8Uw4Lz0b2yXjpOf31awQfXHJthrYB4zppnL3c1 new-key-block -s 850955
```

[link](https://testnet.tonviewer.com/transaction/8b550fa8348188a76d07d94b27ffb080e0009f5d779aa5a94a359e2616626613)

9. send checkTransaction with tx in new epoch -> Ok.

```bash
go run main.go --network testnet --config ./.trustless-bridge-cli.yaml send -a EQCEILr1N8ey9Ar-9OtnCq8A4v217lsE0pJuEZgrZOStAVVa check-tx -s 850956 -t 92C2D323580579391FDE2AB9EE0B749CF6D66D179CE0F0430F049A0735180557
```

[link](https://testnet.tonviewer.com/transaction/8b39492dcc99ec8a9b1833dfc18a9bb610180149dc169fbfcf97aabebb9b538f)