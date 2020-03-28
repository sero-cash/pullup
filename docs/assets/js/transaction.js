var Transaction = {

    currencyDecimal : [],
    currencyDecimalFix : [],

    pkBalance:{},

    loadProperties: function (lang) {
        jQuery.i18n.properties({
            name: 'lang',
            path: 'assets/i18n/',
            mode: 'map',
            language: lang,
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a').text($.i18n.prop('navbar_stake'));
                $('.navbar-nav li:eq(3) a').text($.i18n.prop('navbar_dapps'));

                $('h4:eq(0)').text($.i18n.prop('send_tx'));
                $('label:eq(0)').text($.i18n.prop('send_tx_from'));
                $('label:eq(2)').text($.i18n.prop('send_tx_to'));
                $('label:eq(1)').text($.i18n.prop('send_tx_currency'));
                $('label:eq(3)').text($.i18n.prop('send_tx_amount'));
                $('label:eq(4)').text($.i18n.prop('send_tx_avaliable'));
                $('label:eq(5)').text($.i18n.prop('send_tx_price'));
                $('label:eq(6)').text($.i18n.prop('send_tx_fee'));
                $('label:eq(7)').text($.i18n.prop('send_tx_total'));
                $('#sub1').text($.i18n.prop('send_tx_send'));

                $('.modal-title').text($.i18n.prop('send_tx_titlem'));
                $('.col-lg-3:eq(0)').text($.i18n.prop('send_tx_from'));
                $('.col-lg-3:eq(1)').text($.i18n.prop('send_tx_to'));
                $('.col-lg-3:eq(2)').text($.i18n.prop('send_tx_amount'));
                $('.col-lg-3:eq(3)').text($.i18n.prop('send_tx_fee'));
                $('.col-lg-3:eq(4)').text($.i18n.prop('send_tx_total'));
                $('.col-lg-3:eq(5)').text($.i18n.prop('send_tx_pwd'));

                $('#password').attr('placeholder',$.i18n.prop('send_tx_pwdtips'));
                $('#amount').attr('placeholder',$.i18n.prop('send_tx_amount_tips'));
                $('#address').attr('placeholder',$.i18n.prop('send_tx_address_tips'));

                $('.modal-footer button:eq(0)').text($.i18n.prop('send_tx_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('send_tx_confirm'));

            }
        });
    },


    init: function () {
        var that = this;
        that.getAccountlist();

        $('.close').bind('click', function () {
            $('.modal').hide();
            $("#sub1").attr('disabled', false);
            $('.modal-footer button:eq(1)').unbind('click');
        });


        $('.address').bind('change', function () {
            that.changeAccount();
        });

        $('.currency').bind('change', function () {
            that.changeCurrency();
        });


        $('#amount').bind('input', function () {
            that.calculate();
        });

        $('#gasprice').bind('input', function () {
            that.calculate();
        });

        $('.currency').bind('change', function () {
            that.calculate();
        });

        $('.toast').toast({animation: true, autohide: true, delay: 2000})

    },


    getAccountlist: function () {
        var that = this;
        var biz = {}
        Common.postAsync("account/list", biz, {}, function (res) {
            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var dataArray = res.biz;
                    var _pkBalance = {};
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        if (!$.isEmptyObject(data.Balance)) {
                            var balanceObj = data.Balance;
                            var balanceMap = Object.keys(balanceObj);
                            for (var currency of balanceMap) {
                                if(that.currencyDecimal && that.currencyDecimal[currency]){
                                    continue
                                }
                                if(currency === 'SERO'){
                                    that.currencyDecimal[currency]= new BigNumber(10).pow(new BigNumber(18));
                                    continue
                                }
                                Common.post('decimal', {Currency: currency}, {}, function (res) {
                                    var dcm = res.biz;
                                    if(currency === 'SERO'){ dcm = 18}
                                    that.currencyDecimalFix[currency]= dcm>6?6:dcm;
                                    var decimal = new BigNumber(10).pow(new BigNumber(dcm));
                                    that.currencyDecimal[currency]= decimal;
                                });
                            }
                            _pkBalance[data.PK] = data;
                        }
                    }
                    that.pkBalance = _pkBalance;
                    setTimeout(function () {
                        that.setCurrency();
                    },150);
                }
            }
        });


    },

    setCurrency:function () {
        var that = this;
        var i=0;
        var hasSet = false;
        Object.keys(that.pkBalance).forEach(function(PK){
            var data = that.pkBalance[PK];
            var balanceObj = data.Balance;

            var acName = "Account"+(i + 1);
            if (data.Name){
                acName = data.Name;
            }

            if(!$.isEmptyObject(balanceObj)){
                var balanceMap = Object.keys(balanceObj);
                if (!hasSet){
                    var j=0;
                    for (var currency of balanceMap) {
                        var decimal = that.currencyDecimal[currency];
                        var fix = that.currencyDecimalFix[currency];
                        var balance = new BigNumber(balanceObj[currency]).dividedBy(decimal).toFixed(fix);
                        $('.currency').append(`<option value="${balance}" ${currency === 'SERO'?"selected":""}>${currency} ${balance}</option>`);
                        if(j === 0 ){
                            $('.currencyp').text(balance);
                            $('.currencys').text(currency);
                        }
                        if(currency === 'SERO'){
                            $('.currencyp').text(balance);
                            $('.currencys').text(currency);
                        }
                        j++
                    }
                    hasSet = true
                }
                $('.address').append(`<option value="${PK}" ${currency === 'SERO' ? 'selected' : ''}>${ acName + ": " + PK.substring(0, 20) + ' ... ' + PK.substring(PK.length - 20) }</option>`);
                i++;
            }


        });


    },

    changeAccount: function () {
        var that = this;
        var pk = $(".address").val();

        var biz = {
            PK: pk,
        }

        Common.post("account/detail", biz, {}, function (res) {
            $('.currency').empty();
            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var data = res.biz;
                    var balance = new BigNumber(0).toFixed(6);
                    var hasSet = false;
                    if (data.Balance && !hasSet) {
                        var balanceObj = data.Balance;
                        for (var currency of Object.keys(balanceObj)) {
                            var decimal =  that.currencyDecimal[currency];
                            balance = new BigNumber(balanceObj[currency]).dividedBy(decimal).toFixed(that.currencyDecimalFix[currency]);
                            $('.currency').append(`<option value="${balance}" ${currency === 'SERO' ? 'selected' : ''}>${currency} ${balance}</option>`);
                            if (!hasSet){
                                $('.currencyp').text(balance);
                                $('.currencys').text(currency);
                                hasSet = true;
                            }
                            if (currency === 'SERO'){
                                $('.currencyp').text(balance);
                                $('.currencys').text(currency);
                                hasSet = true;
                            }
                        }
                    }
                }
            }
        })
    },

    changeCurrency: function () {
        var balance = $('.currency').val();
        var currency = $('.currency').find("option:selected").text();
        $('.currencyp').text(balance);
        $('.currencys').text(currency.substring(currency,currency.indexOf(" ")));
    },


    calculate: function () {

        var that = this;
        var amount = $("#amount").val();
        var gasprice = $("#gasprice").val();
        var currency =  $('.currency').find("option:selected").text();
        currency = currency.substring(currency,currency.indexOf(" "));
        var decimal =  that.currencyDecimal[currency];
        if (amount > 0 && gasprice > 0) {
            amount = new BigNumber(amount).multipliedBy(decimal);
            gasprice = new BigNumber(gasprice).multipliedBy(new BigNumber(10).pow(9));
            var fee = gasprice.multipliedBy(25000).dividedBy(Common.baseDecimal);
            var total = amount;
            // $("#amount").val(amount.dividedBy(decimal).toFixed(that.currencyDecimalFix[currency]))
            total = total.dividedBy(decimal);
            if (currency === 'SERO'){
                total = fee.plus(total)
            }
            $('.calculate span:eq(0)').text(fee.toFixed(8));
            $('.calculate span:eq(1)').text(total.toFixed(that.currencyDecimalFix[currency]) + ' ' + currency);
        } else {
            $('.calculate span:eq(0)').text('0.000000');
            $('.calculate span:eq(1)').text('0.000000 ' + currency);
        }
    },


    subTx: function () {
        var that = this;

        var from = $(".address").val();
        var currency = $('.currency').find("option:selected").text();
        currency = currency.substring(currency,currency.indexOf(" "));
        var avliable = $('.currencyp').text();
        avliable = new BigNumber(avliable);
        $('#myModal').modal({backdrop: 'static', keyboard: false});
        var amountStr = $("#amount").val();
        var to = $("#address").val();
        var gasprice = $("#gasprice").val();
        amountStr = new BigNumber(amountStr);
        var decimal =  that.currencyDecimal[currency];

        var amount = amountStr.multipliedBy(decimal);

        gasprice = new BigNumber(gasprice).multipliedBy(new BigNumber(10).pow(9));
        var fee = gasprice.multipliedBy(25000);
        var total = amount;
        if (currency === 'SERO'){
            total = fee.plus(amount);
        }

        $(".modal-body ul li:eq(0) div div:eq(1)").text(from);
        $(".modal-body ul li:eq(1) div div:eq(1)").text(to);
        $(".modal-body ul li:eq(2) div div:eq(1)").text(amount.dividedBy(decimal).toFixed(that.currencyDecimalFix[currency]) + ' ' + currency);
        $(".modal-body ul li:eq(3) div div:eq(1)").text(fee.dividedBy(Common.baseDecimal).toFixed(8) + ' SERO');
        $(".modal-body ul li:eq(4) div div:eq(1)").text(total.dividedBy(decimal).toFixed(that.currencyDecimalFix[currency]) + ' ' + currency);


        $('.modal-footer button:eq(0)').bind('click', function () {
            $('#sub1').attr('disabled', false);
            $('.modal-footer button:eq(1)').unbind('click');
        });

        $('.modal-footer button:eq(1)').bind('click', function () {
            if (total.comparedTo(avliable.multipliedBy(decimal)) > 0) {
                $('.toast:eq(1) div:eq(0)').text($.i18n.prop('send_tx_lessAmount'));
                $('.toast:eq(1)').toast('show')
                $('#sub1').attr('disabled', false);
                $('.modal-footer button:eq(1)').attr('disabled', false);
            } else {
                $('.modal-footer button:eq(1)').attr('disabled', true).text($.i18n.prop('send_tx_sending'));
                var password = $("#password").val();
                if (password === '') {
                    $('.toast:eq(1) div:eq(0)').text($.i18n.prop('send_tx_pwdtips'));
                    $('.toast:eq(1)').toast('show');
                    $('#sub1').attr('disabled', false);
                    $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                } else {
                    var biz = {
                        From: from,
                        To: to,
                        Currency: currency,
                        Amount: amount.toString(10),
                        GasPrice: gasprice.toString(10),
                        Password:password,
                    }
                    Common.postAsync('tx/transfer', biz, {}, function (res) {
                        if (res.base.code === 'SUCCESS') {
                            $('.toast:eq(1) div:eq(0)').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                            $('.toast:eq(1)').toast('show')
                            setTimeout(function () {
                                window.location.href = "account-detail.html?pk=" + from;
                            }, 1500);
                            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                        } else {
                            $('.toast:eq(1) div:eq(0)').text(Common.convertErrors(res.base.desc));
                            $('.toast:eq(1)').toast('show')
                            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                        }
                        $('#sub1').attr('disabled', false);
                    });
                }
            }
        });
    }


}