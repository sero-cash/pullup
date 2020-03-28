var operation = {

    method: {
        init: "init",
        accountDetail: "accountDetail",
        accountList: "accountList",
        executeContract: "executeContract",
        call: "call",
        estimateGas: "estimateGas"
    }

};

var Browser = {

    init: function () {
        var that = this;

        window.addEventListener("message", function () {
            if (event !== undefined && event.data) {
                var msg = event.data;
                console.log("pullup receive msg: ", msg);
                if (msg.method) {
                    if (msg.method === operation.method.init) {
                        msg.data = that.initDApp(msg.data);
                        that.sendMessage(msg);

                    } else if (msg.method === operation.method.accountDetail) {
                        that.getAccountDetail(msg.data, function (data) {
                            msg.data = data;
                            that.sendMessage(msg);
                        });
                    } else if (msg.method === operation.method.accountList) {
                        that.getAccountList(function (data) {
                            msg.data = data;
                            that.sendMessage(msg);
                        });
                    } else if (msg.method === operation.method.executeContract) {
                        that.executeContract(msg.data.tx, function (txHash) {
                            msg.data = txHash;
                            that.sendMessage(msg);
                        });
                    } else if (msg.method === operation.method.call) {
                        that.call(msg.data, function (data) {
                            msg.data = data;
                            that.sendMessage(msg);
                        })
                    } else if (msg.method === operation.method.estimateGas) {
                        that.estimateGas(msg.data, function (data) {
                            msg.data = data;
                            that.sendMessage(msg);
                        })
                    } else {
                        that.sendMessage("operation method is invalid !");
                    }
                } else {
                    that.sendMessage("operation method is required !");
                }
            }
        }, false);

        $('.toast').toast({animation: true, autohide: true, delay: 2000})
    },

    initDApp: function (data) {
        var that = this;
        if (data) {
            // var mainFrame = document.getElementById('ifrModel');
            if (data.name && data.contractAddress && data.github && data.author && data.url && data.logo) {
                if (data.url && data.url.indexOf("http") === -1) {
                    data.url = "http://" + data.url;
                    data.logo = "http://" + data.logo
                }
                that.storageDApp(data)
            }
        }

        return "success"
    },

    getAccountList: function (cb) {
        try {
            Common.post("account/list", {}, {}, function (res) {
                if (res.base.code === 'SUCCESS') {
                    if (cb) {
                        cb(res.biz)
                    }
                }
            });
        } catch (e) {
            alert(e.message);
        }
    },

    getAccountDetail: function (address, cb) {
        try {
            if (address) {
                var biz = {
                    PK: address,
                }
                Common.post("account/detail", biz, {}, function (res) {
                    if (res.base.code === "SUCCESS") {
                        if (cb) {
                            var detail = res.biz;
                            let assetsMap = new Map();
                            if (detail && detail.Balance) {
                                var balanceObj = detail.Balance;
                                for (var k of Object.keys(balanceObj)) {
                                    assetsMap.set(k, balanceObj[k]);
                                }
                                detail.Balance = assetsMap;
                            }
                            cb(detail)
                        }
                    } else {
                        alert(res.base.desc);
                    }
                });
            } else {
                alert("params address undefined");
            }
        } catch (e) {
            alert(e.message);
        }
    },

    call: function (data, cb) {
        try {
            if (data) {
                Common.postSeroRpc("sero_call", [data, "latest"], function (res) {
                    if (cb) {
                        cb(res.result)
                    }
                })
            }
        } catch (e) {
            alert(e.message);
        }
    },

    estimateGas: function (data, cb) {
        try {
            if (data) {
                Common.postSeroRpc("sero_estimateGas", [data], function (res) {
                    if (cb) {
                        cb(res.result)
                    }
                });
            }
        } catch (e) {
            alert(e.message);
        }
    },

    executeContract: function (data, cb) {
        var that = this;
        if (data) {
            var gasPrice = 1000000000;
            if(data.gas_price){
                gasPrice = data.gas_price;
            }else if(data.gasPrice){
                gasPrice = data.gasPrice;
            }
            var fee = new BigNumber(data.gas).multipliedBy(gasPrice).dividedBy(Common.baseDecimal);
            if (data.cy && data.cy !== "SERO") {
                var biz = {
                    Currency: data.cy,
                }
                Common.post('decimal', biz, {}, function (res) {
                    var cDecimal = new BigNumber(10).pow(new BigNumber(res.biz));
                    var amount = new BigNumber(data.value, 16).dividedBy(cDecimal);
                    $('.modal-body ul li:eq(0) div div:eq(1)').text(data.from);
                    $('.modal-body ul li:eq(1) div div:eq(1)').text(data.to);
                    $('.modal-body ul li:eq(2) div div:eq(1)').text(amount + " " + data.cy);
                    $('.modal-body ul li:eq(3) div div:eq(1)').text(data.data);
                    $('.modal-body ul li:eq(4) div div:eq(1)').text(fee + " SERO");
                });
            } else {
                var amount = new BigNumber(data.value, 16).dividedBy(Common.baseDecimal);
                $('.modal-body ul li:eq(0) div div:eq(1)').text(data.from);
                $('.modal-body ul li:eq(1) div div:eq(1)').text(data.to);
                $('.modal-body ul li:eq(2) div div:eq(1)').text(amount + " " + data.cy);
                $('.modal-body ul li:eq(3) div div:eq(1)').text(data.data);
                $('.modal-body ul li:eq(4) div div:eq(1)').text(fee + " SERO");
            }
        }
        $("#transferModal").modal('show');

        $(".modal-footer button:eq(1)").unbind('click').bind('click', function () {
            that.submit(data, cb)
        });
    },

    submit: function (data, cb) {
        $(".modal-footer button:eq(1)").text($.i18n.prop('send_tx_sending')).attr('disabled',true);

        var password = $('#password').val();
        var gasPrice = 1000000000;
        if(data.gas_price){
            gasPrice = data.gas_price;
        }else if(data.gasPrice){
            gasPrice = data.gasPrice;
        }
        if(!password){
            $('.toast div:eq(0)').text($.i18n.prop('send_tx_pwdtips'));
            $('.toast').toast('show');
            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
        }else{
            var biz = {
                From: data.from,
                To: data.to,
                Amount: data.value,
                GasPrice: gasPrice,
                Gas: data.gas,
                Currency: data.cy,
                Data: data.data,
                Password: password,
            }

            Common.postAsync('tx/transfer', biz, {}, function (res) {
                if (res.base.code === 'SUCCESS') {
                    $('.toast div:eq(0)').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                    console.log("res:", res);
                    if (cb) {
                        cb(res.biz);
                    }
                    setTimeout(function () {
                        $("#transferModal").modal('hide');
                    },2000)
                } else {
                    $('.toast div:eq(0)').text(Common.convertErrors(res.base.desc));
                }
                $('.toast').toast('show');
                $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
            });
        }
    },

    storageDApp: function (data) {
        var that = this;
        try {
            var dappListKey = "dapp_list";
            var list = localStorage.getItem(dappListKey);
            if (!list || list.length === 0) {
                list = [data.contractAddress]
                localStorage.setItem(dappListKey, list)
            } else {
                var tempList = [];
                var has = false;
                for (var v of list) {
                    if (v === data.contractAddress) {
                        has = true;
                        break
                    }
                }
                if (!has) {
                    tempList.push(data.contractAddress)
                }
                localStorage.setItem(dappListKey, tempList.concat(list))
            }

            localStorage.setItem("dapp:" + data.name + data.contractAddress, data)
        } catch (e) {
            alert(e.message);
        }
    },

    sendMessage: function (msg) {
        console.log("pullup send msg: ", msg);
        var childFrameObj = document.getElementById('myFrame');
        childFrameObj.contentWindow.postMessage(msg, '*');
    },

    loadProperties: function (lang) {
        jQuery.i18n.properties({
            name: 'lang',
            path: 'assets/i18n/',
            mode: 'map',
            language: lang,
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a span').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a span').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a span').text($.i18n.prop('navbar_stake'));
                $('.navbar-nav li:eq(3) a span').text($.i18n.prop('navbar_dapps'));

                $('.modal-title').text($.i18n.prop('send_tx_titlem'));
                $('.col-lg-3:eq(0)').text($.i18n.prop('send_tx_from'));
                $('.col-lg-3:eq(1)').text($.i18n.prop('send_tx_to'));
                $('.col-lg-3:eq(2)').text($.i18n.prop('send_tx_amount'));
                $('.col-lg-3:eq(4)').text($.i18n.prop('send_tx_fee'));
                $('#password').attr('placeholder',$.i18n.prop('send_tx_pwdtips'));
                $('.modal-footer button:eq(0)').text($.i18n.prop('send_tx_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('send_tx_confirm'));

            }
        });
    },

};