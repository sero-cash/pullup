var Token = {

    solidity: `pragma solidity ^0.4.16;
        import "./seroInterface.sol";
        
        
        /**
         * Math operations with safety checks
         */
        contract SafeMath {
            function safeMul(uint256 a, uint256 b) pure internal returns (uint256) {
                uint256 c = a * b;
                assert(a == 0 || c / a == b);
                return c;
            }
        
            function safeDiv(uint256 a, uint256 b) pure internal returns (uint256) {
                assert(b > 0);
                uint256 c = a / b;
                assert(a == b * c + a % b);
                return c;
            }
        
            function safeSub(uint256 a, uint256 b) pure internal returns (uint256) {
                assert(b <= a);
                return a - b;
            }
        
            function safeAdd(uint256 a, uint256 b) pure internal returns (uint256) {
                uint256 c = a + b;
                assert(c>=a && c>=b);
                return c;
            }
        }
        
        
        contract owned {
            address public owner;
        
            constructor() public {
                owner = msg.sender;
            }
        
            modifier onlyOwner {
                require(msg.sender == owner);
                _;
            }
        
            function transferOwnership(address newOwner) onlyOwner public {
                owner = newOwner;
            }
        }
        
        
        contract SOSOSO is SeroInterface ,owned ,SafeMath { 
            
            string private _name;
            string private _symbol;
            uint8 private _decimals;
           uint256 private _totalSupply;
        
        
            
            mapping (address => mapping (address => uint256)) private allowance;
        
            // This generates a public event on the blockchain that will notify clients
            event Transfer(address indexed from, address indexed to, uint256 value);
            
            // This generates a public event on the blockchain that will notify clients
            event Approval(address indexed _owner, address indexed _spender, uint256 _value);
        
            /**
             * Constrctor function
             *
             * Initializes contract with initial supply tokens to the creator of the contract
             */
            constructor(
                uint256 initialSupply,
                string tokenName,
                string tokenSymbol,
                uint8 decimals
            ) public payable{
                _totalSupply = initialSupply * 10 ** uint256(decimals);
                require(sero_issueToken(_totalSupply,tokenSymbol));
                _name = tokenName;                                       // Set the name for display purposes
                _symbol = tokenSymbol;                               // Set the currency for display purposes
                _decimals = decimals;
                
               
            }
            
            /**
             * @return the name of the token.
             */
            function name() public view returns (string memory) {
                return _name;
            }
        
            /**
             * @return the symbol of the token.
             */
            function symbol() public view returns (string memory) {
                return _symbol;
            }
        
            /**
             * @return the number of decimals of the token.
             */
            function decimals() public view returns (uint8) {
                return _decimals;
            }
            
            function totalSupply() public view returns (uint256) {
                return _totalSupply;
            }
            
        
            /**
             * the contract current left balance 
             */
            function balanceOf() public returns(uint256 amount) {
                return sero_balanceOf(_symbol);
            }
            
            /**
             * Transfer tokens
             *
             * Send \`_value\` tokens to \`_to\` from your account
             *
             * @param _to The address of the recipient
             * @param _value the amount to send
             */
            function transfer(address _to, uint256 _value) public onlyOwner returns (bool success) {
                return sero_send(_to,_symbol,_value,'','');
            }
        
            /**
             * Transfer tokens from other address
             *
             * Send \`_value\` tokens to \`_to\` in behalf of \`_from\`
             *
             * @param _from The address of the sender
             * @param _to The address of the recipient
             * @param _value the amount to send
             */
            function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
                require(_value <= allowance[_from][msg.sender]);     // Check allowance
                require (sero_send(_to,_symbol,_value,'',''));
                allowance[_from][msg.sender] -= _value;
                return true ;
            }
        
            /**
             * Set allowance for other address
             *
             * Allows \`_spender\` to spend no more than \`_value\` tokens in your behalf
             *
             * @param _spender The address authorized to spend
             * @param _value the max amount they can spend
             */
            function approve(address _spender, uint256 _value) public onlyOwner
                returns (bool success) {
                allowance[msg.sender][_spender] =_value;
                emit Approval(msg.sender, _spender, _value);
                return true;
            }
            
            
            function withDraw(address _to,string _cy) public onlyOwner{
                uint256 balance = sero_balanceOf(_cy);
                require(sero_balanceOf(_cy)> 0);
                require(sero_send_token(_to,_cy,balance));
            }
        }
    `,
    abi: [{
        "constant": true,
        "inputs": [],
        "name": "name",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "_spender", "type": "address"}, {"name": "_value", "type": "uint256"}],
        "name": "approve",
        "outputs": [{"name": "success", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [],
        "name": "totalSupply",
        "outputs": [{"name": "", "type": "uint256"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "_from", "type": "address"}, {"name": "_to", "type": "address"}, {
            "name": "_value",
            "type": "uint256"
        }],
        "name": "transferFrom",
        "outputs": [{"name": "success", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [],
        "name": "decimals",
        "outputs": [{"name": "", "type": "uint8"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "_to", "type": "address"}, {"name": "_cy", "type": "string"}],
        "name": "withDraw",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [],
        "name": "balanceOf",
        "outputs": [{"name": "amount", "type": "uint256"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [],
        "name": "owner",
        "outputs": [{"name": "", "type": "address"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [{"name": "x", "type": "bytes32"}],
        "name": "bytes32ToString",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "pure",
        "type": "function"
    }, {
        "constant": true,
        "inputs": [],
        "name": "symbol",
        "outputs": [{"name": "", "type": "string"}],
        "payable": false,
        "stateMutability": "view",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "_to", "type": "address"}, {"name": "_value", "type": "uint256"}],
        "name": "transfer",
        "outputs": [{"name": "success", "type": "bool"}],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "constant": false,
        "inputs": [{"name": "newOwner", "type": "address"}],
        "name": "transferOwnership",
        "outputs": [],
        "payable": false,
        "stateMutability": "nonpayable",
        "type": "function"
    }, {
        "inputs": [{"name": "initialSupply", "type": "uint256"}, {
            "name": "tokenName",
            "type": "string"
        }, {"name": "tokenSymbol", "type": "string"}, {"name": "decimals", "type": "uint8"}],
        "payable": true,
        "stateMutability": "payable",
        "type": "constructor"
    }, {
        "anonymous": false,
        "inputs": [{"indexed": true, "name": "from", "type": "address"}, {
            "indexed": true,
            "name": "to",
            "type": "address"
        }, {"indexed": false, "name": "value", "type": "uint256"}],
        "name": "Transfer",
        "type": "event"
    }, {
        "anonymous": false,
        "inputs": [{"indexed": true, "name": "_owner", "type": "address"}, {
            "indexed": true,
            "name": "_spender",
            "type": "address"
        }, {"indexed": false, "name": "_value", "type": "uint256"}],
        "name": "Approval",
        "type": "event"
    }],
    data: '0x608060408190527f3be6bf24d822bcd6f6348f6f5a5c2d3108f04991ee63e80cde49a8c4746a0ef36000557fcf19eb4256453a4e30b6a06d651f1970c223fb6bd1826a28ed861f0e602db9b86001557f868bd6629e7c2e3d2ccf7b9968fad79b448e7a2bfb3ee20ed1acbc695c3c8b236002557f7c98e64bd943448b4e24ef8c2cdec7b8b1275970cfe10daf2a9bfa4b04dce9056003557fa6a366f1a72e1aef5d8d52ee240a476f619d15be7bc62d3df37496025b83459f6004557ff1964f6690a0536daa42e5c575091297d2479edcc96f721ad85b95358644d2766005557f9ab0d7c07029f006485cf3468ce7811aa8743b5a108599f6bec9367c50ac6aad6006557fa6cafc6282f61eff9032603a017e652f68410d3d3c69f0a3eeca8f181aec1d176007557f6800e94e36131c049eaeb631e4530829b0d3d20d5b637c8015a8dc9cedd70aed6008557fbbf1aa2159b035802d0a4d44611849d5d4ada0329c81580477d5ec3e82f4f0a66009557fa8b83585a613dcf6c905ad7e0ce34cd07d1283cc72906d1fe78037d49adae455600a55610d573881900390819083398101604090815281516020830151918301516060840151600b8054600160a060020a0319163317905560ff8116600a0a8302600f81905592949384019391909101916101f39083640100000000610243810204565b15156101fe57600080fd5b825161021190600c906020860190610283565b50815161022590600d906020850190610283565b50600e805460ff191660ff929092169190911790555061031e915050565b60408051818152606080820183526000929091906020820161080080388339019050509050828152836020820152600054604082a1602001519392505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102c457805160ff19168380011785556102f1565b828001600101855582156102f1579182015b828111156102f15782518255916020019190600101906102d6565b506102fd929150610301565b5090565b61031b91905b808211156102fd5760008155600101610307565b90565b610a2a8061032d6000396000f3006080604052600436106100b95763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166306fdde0381146100be578063095ea7b31461014857806318160ddd1461018057806323b872dd146101a7578063313ce567146101d15780635413e8a8146101fc578063722713f7146102655780638da5cb5b1461027a5780639201de55146102ab57806395d89b41146102c3578063a9059cbb146102d8578063f2fde38b146102fc575b600080fd5b3480156100ca57600080fd5b506100d361031d565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561010d5781810151838201526020016100f5565b50505050905090810190601f16801561013a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561015457600080fd5b5061016c600160a060020a03600435166024356103b3565b604080519115158252519081900360200190f35b34801561018c57600080fd5b50610195610434565b60408051918252519081900360200190f35b3480156101b357600080fd5b5061016c600160a060020a036004358116906024351660443561043a565b3480156101dd57600080fd5b506101e6610555565b6040805160ff9092168252519081900360200190f35b34801561020857600080fd5b5060408051602060046024803582810135601f8101859004850286018501909652858552610263958335600160a060020a031695369560449491939091019190819084018382808284375094975061055e9650505050505050565b005b34801561027157600080fd5b506101956105b3565b34801561028657600080fd5b5061028f610651565b60408051600160a060020a039092168252519081900360200190f35b3480156102b757600080fd5b506100d3600435610660565b3480156102cf57600080fd5b506100d3610819565b3480156102e457600080fd5b5061016c600160a060020a036004351660243561087a565b34801561030857600080fd5b50610263600160a060020a03600435166108fe565b600c8054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103a95780601f1061037e576101008083540402835291602001916103a9565b820191906000526020600020905b81548152906001019060200180831161038c57829003601f168201915b5050505050905090565b600b54600090600160a060020a031633146103cd57600080fd5b336000818152601060209081526040808320600160a060020a03881680855290835292819020869055805186815290519293927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925929181900390910190a350600192915050565b600f5490565b600160a060020a038316600090815260106020908152604080832033845290915281205482111561046a57600080fd5b600d805460408051602060026001851615610100026000190190941693909304601f8101849004840282018401909252818152610516938793919290918301828280156104f85780601f106104cd576101008083540402835291602001916104f8565b820191906000526020600020905b8154815290600101906020018083116104db57829003601f168201915b50505050508460206040519081016040528060008152506000610944565b151561052157600080fd5b50600160a060020a038316600090815260106020908152604080832033845290915290208054829003905560019392505050565b600e5460ff1690565b600b54600090600160a060020a0316331461057857600080fd5b6105818261099c565b9050600061058e8361099c565b1161059857600080fd5b6105a38383836109d3565b15156105ae57600080fd5b505050565b600d805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815260009361064c93919290918301828280156106425780601f1061061757610100808354040283529160200191610642565b820191906000526020600020905b81548152906001019060200180831161062557829003601f168201915b505050505061099c565b905090565b600b54600160a060020a031681565b60408051602080825281830190925260609160009183918391829184919080820161040080388339019050509350600092505b602083101561072f576008830260020a870291507fff00000000000000000000000000000000000000000000000000000000000000821615610719578184868151811015156106de57fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600190940193610724565b84156107245761072f565b600190920191610693565b846040519080825280601f01601f19166020018201604052801561075d578160200160208202803883390190505b509050600092505b8483101561080f57838381518110151561077b57fe5b9060200101517f010000000000000000000000000000000000000000000000000000000000000090047f01000000000000000000000000000000000000000000000000000000000000000281848151811015156107d457fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350600190920191610765565b9695505050505050565b600d8054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103a95780601f1061037e576101008083540402835291602001916103a9565b600b54600090600160a060020a0316331461089457600080fd5b600d805460408051602060026001851615610100026000190190941693909304601f81018490048402820184019092528181526108f7938793919290918301828280156104f85780601f106104cd576101008083540402835291602001916104f8565b9392505050565b600b54600160a060020a0316331461091557600080fd5b600b805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a0392909216919091179055565b6040805160a080825260c0820190925260009160609190602082016114008038833901905050905086815285602082015284604082015283606082015282608082015260025460a082a1608001519695505050505050565b6040805160208082528183019092526000916060919080820161040080388339019050509050828152600154602082a15192915050565b60006109f684848460206040519081016040528060008152506000600102610944565b9493505050505600a165627a7a72305820a13a2a4571c355bbc6b1ba270d1d2ee11d4703ace504af59c11cad9e891d37d70029',

    coinFee: 0,
    gas: 0,
    account: [],

    init: function () {
        var that = this;
        that.getAccountlist();
        that.getTokenList();

        $('#tokenSymbol').bind('input', function () {
            var symbol = $(this).val();
            if (symbol.length === 4) {
                that.coinFee = 1000;
                $('.coinFee').text("1000.00");
            } else if (symbol.length === 5) {
                that.coinFee = 10;
                $('.coinFee').text("10.00");
            } else if (symbol.length === 6) {
                that.coinFee = 1;
                $('.coinFee').text("1.00");
            } else if (symbol.length >= 7) {
                that.coinFee = 0.1;
                $('.coinFee').text("0.10");
            } else {
                that.coinFee = 0;
                $('.coinFee').text("0.00");
            }
            that.estimateGas();
        });
        $('#tokenName').bind('input', function () {
            that.estimateGas();
        });
        $('#totalSupply').bind('input', function () {
            that.estimateGas();
        });
        $('#tokenDecimals').bind('input', function () {
            that.estimateGas();
        });

        $('.issue').bind('click', function () {
            $('.modal-deploy').modal('show');
        });

        $('.watch').bind('click', function () {
            $('.modal-watch').modal('show');
        });

        $('#contractAddress').bind('input', function () {
            that.loadContract();
        });


        $('#toast1').toast({animation: true, autohide: true, delay: 2000})
        $('#toast2').toast({animation: true, autohide: true, delay: 2000})
        $('#toast3').toast({animation: true, autohide: true, delay: 2000})
    },

    loadProperties: function (lang) {

        jQuery.i18n.properties({
            name: 'lang', // 资源文件名称
            path: '../../assets/i18n/', // 资源文件所在目录路径
            mode: 'map', // 模式：变量或 Map
            language: lang, // 对应的语言
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a').text($.i18n.prop('navbar_stake'));
                $('.navbar-nav li:eq(3) a').text($.i18n.prop('navbar_dapps'));

                $('ol li:eq(0) a').text($.i18n.prop('navbar_dapps'));
                $('ol li:eq(1)').text($.i18n.prop('dapp_token_bread'));

                $('.btn-issue').text($.i18n.prop('dapp_token_button_issue'));
                $('.btn-watch').text($.i18n.prop('dapp_token_button_watch'));

                $('thead tr td:eq(0)').text($.i18n.prop('dapp_token_table_address'));
                $('thead tr td:eq(1)').text($.i18n.prop('dapp_token_table_name'));
                $('thead tr td:eq(2)').text($.i18n.prop('dapp_token_table_symbol'));
                $('thead tr td:eq(3)').text($.i18n.prop('dapp_token_table_decimal'));
                $('thead tr td:eq(4)').text($.i18n.prop('dapp_token_table_total'));
                $('thead tr td:eq(5)').text($.i18n.prop('dapp_token_table_balance'));
                $('thead tr td:eq(6)').text($.i18n.prop('dapp_token_table_account'));
                $('thead tr td:eq(7)').text($.i18n.prop('dapp_token_table_operation'));

                $('.modal-title:eq(0)').text($.i18n.prop('dapp_token_modal_issue_title'));

                $('#myModal label:eq(0)').text($.i18n.prop('dapp_token_modal_issue_from'));
                $('#myModal label:eq(1)').text($.i18n.prop('dapp_token_modal_issue_name'));
                $('#myModal label:eq(2)').text($.i18n.prop('dapp_token_modal_issue_symbol'));
                $('#myModal label:eq(3)').text($.i18n.prop('dapp_token_modal_issue_total'));
                $('#myModal label:eq(4)').text($.i18n.prop('dapp_token_modal_issue_decimal'));
                // $('#myModal label:eq(5)').text($.i18n.prop('dapp_token_modal_issue_password'));

                $('.coin-fee').text($.i18n.prop('dapp_token_modal_issue_coinfee'));
                $('.gas-fee').text($.i18n.prop('dapp_token_modal_issue_gasfee'));
                $('.is-total').text($.i18n.prop('dapp_token_modal_issue_total'));
                $('.is-tips').text($.i18n.prop('dapp_token_modal_issue_tips'));

                $('.modal-title:eq(1)').text($.i18n.prop('dapp_token_modal_watch_title'));
                $('.ct-addr').text($.i18n.prop('dapp_token_modal_watch_address'));
                $('.w-st small:eq(0)').text($.i18n.prop('dapp_token_modal_watch_name'));
                $('.w-st small:eq(2)').text($.i18n.prop('dapp_token_modal_watch_symbol'));
                $('.w-st small:eq(4)').text($.i18n.prop('dapp_token_modal_watch_decimal'));
                $('.w-st small:eq(6)').text($.i18n.prop('dapp_token_modal_watch_total'));
                $('.w-st small:eq(8)').text($.i18n.prop('dapp_token_modal_watch_balance'));

                $('.modal-title:eq(2)').text($.i18n.prop('dapp_token_modal_transfer_title'));
                $('#myModal3 label:eq(0)').text($.i18n.prop('dapp_token_modal_transfer_address'));
                $('#myModal3 label:eq(1)').text($.i18n.prop('dapp_token_modal_transfer_name'));
                $('#myModal3 label:eq(2)').text($.i18n.prop('dapp_token_modal_transfer_symbol'));
                $('#myModal3 label:eq(3)').text($.i18n.prop('dapp_token_modal_transfer_total'));

            }
        });
    },

    getAccountlist: function () {

        var that = this;
        var biz = {}
        Common.postAsync("account/list", biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var dataArray = res.biz;
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        if (that.account.length === 0) {
                            that.account.push(data);
                        }
                        var balance = new BigNumber(0).toFixed(6);
                        var acName = "Account" + (i + 1);
                        if (data.Name) {
                            acName = data.Name;
                        }
                        if (data.Balance) {
                            var balanceObj = data.Balance;
                            for (var currency of Object.keys(balanceObj)) {
                                if (currency === 'SERO') {
                                    balance = new BigNumber(balanceObj[currency]).dividedBy(Common.baseDecimal).toFixed(6);
                                    $('#from').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${acName + ": " + data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8)} -- ${balance + ' ' + currency}</option>`);
                                    $('#e_c_from').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${acName + ": " + data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8)} -- ${balance + ' ' + currency}</option>`);
                                }
                            }
                        }
                    }
                }
            }
        })
    },

    getTokenList: function () {

        var that = this;
        var biz = {}
        Common.postPullupRpc("get_tokens", biz, function (res) {
            if (res.result && res.result.length > 0) {
                $('tbody').empty();
                var tokens = res.result;

                for (let token of tokens) {
                    if (token.ContractAddress) {
                        var hidden = false;
                        that.execute(token.ContractAddress, 'balanceOf', [], function (res) {
                            if (res.result) {
                                var tokenBlance = res.result;
                                Common.postSeroRpc("sero_getBalance",[token.ContractAddress,"latest"],function (res) {
                                    var seroS=0;
                                    if(res.result.tkn){
                                        seroS=res.result.tkn["SERO"];
                                    }
                                    $('tbody').append(`
                                        <tr>
                                        <td class="text-break">${token.ContractAddress}</td>
                                        <td>${token.Name}</td>
                                        <td>${token.Symbol}</td>
                                        <td>${token.Decimal}</td>
                                        <td>${new BigNumber(token.Total, 16).toFixed(0)}</td>
                                        <td>${new BigNumber(tokenBlance).dividedBy(new BigNumber(10).pow(parseInt(token.Decimal))).toFixed(6)}</td>
                                        <td>${!!hidden?"unknown":new BigNumber(seroS,16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                        <td><button class="btn btn-outline-info" onclick="showTokenModal(${"'" + token.ContractAddress + "'," + token.Decimal})">Transfer</button></td>
                                        </tr>
                                    `);
                                })
                            }else {
                                Common.postSeroRpc("sero_getBalance",[token.ContractAddress,"latest"],function (res) {
                                    var seroS=0;
                                    if(res.result.tkn){
                                        seroS=res.result.tkn["SERO"];
                                    }
                                    $('tbody').append(`
                                        <tr>
                                        <td class="text-break">${token.ContractAddress}</td>
                                        <td>${token.Name}</td>
                                        <td>${token.Symbol}</td>
                                        <td>${token.Decimal}</td>
                                        <td>${new BigNumber(token.Total, 16).toFixed(0)}</td>
                                        <td>0.000000</td>
                                        <td>${!!hidden?"unknown":new BigNumber(seroS,16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                        <td></td>
                                        </tr>
                                    `);
                                })
                            }
                        });


                    }
                }
            }
        });
    },

    deploy: function () {
        var that = this;
        var from = $('#from').find('option:selected').val();
        var _totalSupply = $('#totalSupply').val();
        var _tokenName = $('#tokenName').val();
        var _tokenSymbol = $('#tokenSymbol').val();
        var _tokenDecimals = $('#tokenDecimals').val();
        if (_totalSupply && _tokenName && _tokenSymbol && _tokenDecimals) {
            var totalSupply = "0x" + new BigNumber(_totalSupply).toString(16);
            var tokenDecimals = "0x" + new BigNumber(_tokenDecimals).toString(16);
            var params = [
                that.abi,
                that.data
            ]
            var args = [];
            args.push(totalSupply, _tokenName.toUpperCase(), _tokenSymbol.toUpperCase(), tokenDecimals);
            params.push(args);
            Common.postSeroRpcSync("sero_packConstruct", params, function (res) {
                var data = res.result;
                if (res.result) {
                    var param = {
                        from: from,
                        data: data,
                        value: "0x" + new BigNumber($('.coinFee').text()).multipliedBy(Common.baseDecimal).toString(16),
                    }
                    Common.postSeroRpcSync("sero_estimateGas", [param], function (res) {

                        if (res.result) {
                            // var password = $('#password').val();
                            var contract_tx_req = {
                                from: from,
                                value: new BigNumber(that.coinFee).multipliedBy(Common.baseDecimal).toString(10),
                                gas_price: "1000000000",
                                gas: new BigNumber(res.result).toString(10),
                                data: data,
                                token: {
                                    Name: _tokenName,
                                    Symbol: _tokenSymbol,
                                    Decimal: _tokenDecimals,
                                    Total: new BigNumber(_totalSupply).multipliedBy(new BigNumber(10).pow(_tokenDecimals).toString(10)),
                                }
                            }

                            var biz = {
                                // password: password,
                                contract_tx_req: contract_tx_req,
                            }

                            var i = 1;
                            var interid = setInterval(function () {
                                $('#sub1').text(`SENDING ${i++}S`);
                            }, 1000)
                            Common.postPullupRpc("deploy_contract", biz, function (res) {
                                if (res.result) {
                                    $('#toast1 div:eq(0)').removeClass('alert-danger').addClass('alert-success').text("successful");
                                    $('.toast .deploy').toast('show');
                                    $('#sub1').text('Confirm');
                                    $('#sub1').attr('disabled', false);
                                    setTimeout(function () {
                                        window.location.href = "../../account-detail.html?pk=" + from;
                                    }, 1500);
                                    clearInterval(interid)
                                } else if (res.error) {
                                    $('#toast1 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text(res.error.message);
                                    $('.toast .deploy').toast('show');
                                    $('#sub1').text('Confirm');
                                    $('#sub1').attr('disabled', false);
                                    clearInterval(interid)
                                }
                            });
                        }
                    });
                }
            });
        }
    },

    estimateGas: function () {
        var that = this;
        var _totalSupply = $('#totalSupply').val();
        var _tokenName = $('#tokenName').val();
        var _tokenSymbol = $('#tokenSymbol').val();
        var _tokenDecimals = $('#tokenDecimals').val();
        var from = $('#from').find('option:selected').val();
        $('.gasFee').text(0.00);
        $('.total').text(0.00);
        if (_tokenName && _tokenSymbol && _tokenDecimals) {
            var totalSupply = "0x" + new BigNumber(_totalSupply).toString(16);
            var tokenDecimals = "0x" + new BigNumber(_tokenDecimals).toString(16);
            var params = [
                that.abi,
                that.data
            ]
            var args = [];
            args.push(totalSupply, _tokenName.toUpperCase(), _tokenSymbol.toUpperCase(), tokenDecimals);
            params.push(args);
            Common.postSeroRpc("sero_packConstruct", params, function (res) {
                if (res.result) {
                    var param = {
                        from: from,
                        data: res.result,
                        value: "0x" + new BigNumber($('.coinFee').text()).multipliedBy(Common.baseDecimal).toString(16),
                    }
                    Common.postSeroRpc("sero_estimateGas", [param], function (res) {
                        if (res.result) {
                            that.gas = res.result;
                            $('.gas').text(res.result);
                            var gas = new BigNumber(res.result);
                            var gasPrice = new BigNumber(1000000000);
                            $('.gasFee').text(gas.multipliedBy(gasPrice).dividedBy(Common.baseDecimal).toFixed(6));
                            $('.total').text(gas.multipliedBy(gasPrice).plus(new BigNumber($('.coinFee').text()).multipliedBy(Common.baseDecimal)).dividedBy(Common.baseDecimal).toFixed(6));
                        } else {
                            $('#toast1 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text("Parameters error or Token symbol already exists!");
                            $('#toast1').toast('show');
                        }
                    });
                }
            });
        }
    },

    isToken: true,
    t_name: '',
    t_symbol: '',
    t_total: 0,
    t_decimals: 0,


    loadContract: function () {
        var that = this;
        var _contractAddress = $('#contractAddress').val();
        if (_contractAddress) {
            $('.w_t_decimal').text("0");
            $('.w_t_symbol').text("");
            $('.w_t_name').text("");
            $('.w_t_total').text("0");
            $('.w_t_balance').text("0");

            that.t_name= '';
            that.t_symbol= '';
            that.t_total= 0;
            that.t_decimals= 0;


            that.execute(_contractAddress, 'decimals', [], function (res) {
                if (res.result) {
                    var decimal = new BigNumber(10).pow(new BigNumber(res.result, 10));
                    that.t_decimals = new BigNumber(res.result, 10);
                    $('.w_t_decimal').text(new BigNumber(res.result, 10).toString(10));
                    that.execute(_contractAddress, 'symbol', [], function (res) {
                        if (res.result) {
                            $('.w_t_symbol').text(res.result);
                            that.t_symbol = res.result;
                        } else {
                            that.isToken = false;
                        }
                    });

                    that.execute(_contractAddress, 'name', [], function (res) {
                        if (res.result) {
                            $('.w_t_name').text(res.result);
                            that.t_name = res.result;
                        } else {
                            that.isToken = false;
                        }
                    });

                    that.execute(_contractAddress, 'totalSupply', [], function (res) {
                        if (res.result) {
                            console.log("res totalSupply: ",res);
                            $('.w_t_total').text(new BigNumber(res.result, 10).dividedBy(decimal).toFixed(6));
                            that.t_total = "0x" + new BigNumber(res.result, 10).dividedBy(decimal).toString(16);
                        } else {
                            that.isToken = false;
                        }
                    });

                    that.execute(_contractAddress, 'balanceOf', [], function (res) {
                        if (res.result) {
                            $('.w_t_balance').text(new BigNumber(res.result, 10).dividedBy(decimal).toFixed(6));

                        } else {
                            that.isToken = false;
                        }
                    });

                } else {
                    Common.postSeroRpc("sero_getBalance",[_contractAddress,"latest"],function (res) {
                        var seroS=0;
                        if(res.result.tkn){
                            seroS=res.result.tkn["SERO"];
                        }

                        that.t_name= 'None';
                        that.t_symbol= 'None';
                        that.t_total= 0;
                        that.t_decimals= 0;

                        $('.account_sero').text(new BigNumber(seroS,16).dividedBy(Common.baseDecimal).toFixed(6));
                    })

                    that.isToken = false;
                }
            });
        } else {
            that.isToken = false;
        }
    },

    addToken: function () {
        var that = this;
        var _contractAddress = $('#contractAddress').val();

        var param = {
            ContractAddress: _contractAddress,
            Name: that.t_name[0],
            Symbol: that.t_symbol[0],
            Decimal: parseInt(that.t_decimals),
            Total: that.t_total===0?"0x0":that.t_total,
        }
        Common.postPullupRpc("watch_tokens", param, function (res) {
            if (res.result) {
                $('#toast2 div:eq(0)').removeClass('alert-danger').addClass('alert-success').text("Add token Successful");
                $('#toast2').toast('show');
                $('#sub2').attr("disabled",false);
                setTimeout(function () {
                    $("#myModa2").modal("hide");
                },2000);
                that.getTokenList();
            }
            if (res.error) {
                $('#toast2 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text(res.error.message);
                $('#toast2').toast('show');
                $('#sub2').attr("disabled",false);
            }
        })

        // if(that.t_name){
        //
        // }else{
        //     $('#toast2 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text("It is not TOKEN contract address");
        //     $('#toast2').toast('show');
        //     $('#sub2').attr("disabled",false);
        // }

    },

    currentContractAddress: '',
    currentDecimal: 0,

    showTransferModal: function (contractAddress, decimal) {
        var that = this;
        that.currentContractAddres = contractAddress;
        that.currentDecimal = decimal;
        var params = [
            "2V1GBFM8szz41kqLHf6CpxG6oBQm2oGFy13YKnPU4cJiSSagbovL8A48Yewkr5x2BD3zcpSfxdKoZm96i3LvFitLuw47iu7KcqwLsEYrDuTKuFPrZiSNQ4hkpzKdz7ZNTaGf",
            "0x1"
        ];
        var param = [
            that.abi,
            contractAddress,
            "transfer",
            params
        ]
        Common.postSeroRpcSync("sero_packMethod", param, function (res) {

            if(res.result){
                $('.modal-execute').modal('show');
            }else{
                alert("No transfer interface available");
            }
        });
    },

    transfer: function () {
        var that = this;
        var _contractAddress = that.currentContractAddres;
        var _toAddress = $('#e_c_to').val();
        var _value = $('#e_c_amount').val();
        // var _password = $('#e_c_password').val();
        console.log("that.currentDecimal: ", that.currentDecimal, _value);
        var _decimal = new BigNumber(10).pow(that.currentDecimal);
        var from = $('#e_c_from').find('option:selected').val();
        if (_contractAddress && _toAddress && _value && from) {
            var params = [
                _toAddress,
                "0x" + new BigNumber(_value).multipliedBy(_decimal).toString(16)
            ];

            var mainPkr = '';
            for (let ac of that.account) {
                if (ac.PK === from) {
                    mainPkr = ac.MainPKr;
                    break;
                }
            }
            if (mainPkr) {
                that.execute(_contractAddress, 'transfer', params, function (res) {
                    if (res.result) {
                        var data = res.result;
                        var executeData = {
                            contract_tx_req: {
                                from: from,
                                to: _contractAddress,
                                value: "0",
                                data: data,
                                gas_price: "1000000000",
                            },
                        };
                        var estimateParam = {
                            from: mainPkr,
                            to: _contractAddress,
                            data: data,
                            value: "0x0",
                        }
                        Common.postSeroRpc("sero_estimateGas", [estimateParam], function (res) {
                            if (res.result) {
                                executeData.contract_tx_req["gas"] = new BigNumber(res.result, 16).toString(10);
                                Common.postPullupRpc("execute_contract", executeData, function (res) {
                                    $('#sub3').attr('disabled', false);
                                    console.log(res);
                                    if (res.result) {
                                        var txHash = res.result;
                                        $('#toast3 div:eq(0)').removeClass('alert-danger').addClass('alert-success').text(txHash);
                                        $('#toast3').toast('show');
                                    }
                                    if (res.error) {
                                        $('#toast3 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text(res.error.message);
                                        $('#toast3').toast('show');
                                    }
                                })
                            }else{
                                $('#sub3').attr('disabled', false);
                                $('#toast3 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('dapp_token_modal_transfer_err'));
                                $('#toast3').toast('show');
                            }
                        });

                    } else {
                        //not token contract
                        $('#sub3').attr('disabled', false);
                        $('#toast3 div:eq(0)').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('dapp_token_modal_transfer_err'));
                        $('#toast3').toast('show');
                    }
                });
            }
        }
    },

    execute: function (contractAddress, methodName, args, callback) {
        var that = this;
        var param = [
            that.abi,
            contractAddress,
            methodName,
            args
        ]
        Common.postSeroRpcSync("sero_packMethod", param, function (res) {
            if (res.result) {
                if (methodName === 'transfer') {
                    if (callback) {
                        callback(res);
                    }
                } else {
                    var callParams = {
                        from: that.account[0].PK,
                        to: contractAddress,
                        data: res.result
                    }
                    Common.postSeroRpcSync("sero_call", [callParams, "latest"], function (res) {
                        if (res) {
                            that.unPackData(contractAddress, methodName, res.result, function (res) {
                                if (callback) {
                                    callback(res);
                                }
                            })
                        }
                    });
                }
            }else{
                if (callback) {
                    callback(res);
                }
            }
        });
    }
    ,

    unPackData: function (contractAddress, methodName, data, callback) {
        var that = this;
        var param = [
            that.abi,
            methodName,
            data
        ]
        Common.postSeroRpcSync("sero_unPack", param, function (res) {
            if (res) {
                callback(res)
            }
        });
    }
    ,

}
