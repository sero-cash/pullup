var StakeHome = {

    init: function () {
        var that = this;

        that.stakeList();
        that.getAccountlist();

        $('.close').bind('click', function () {
            $('.modal').hide();
            $('.modal-footer button:eq(1)').unbind('click');
            $("#sub1").attr('disabled', false);
        });

        $('.register').bind('click', function () {
            window.location.href = 'stake-register.html';
        });

        $('.showShareDetail').bind('click', function () {
            window.location.href = 'stake-detail.html';
        });

        setTimeout(function () {
            that.getAccountlist();
        },10000);

        setTimeout(function () {
            $('.buyShare').bind('click', function () {
                var poolId = $(this).attr('attpoolid');
                window.location.href = 'stake-buy.html?id=' + poolId;
            });
        },100)

    },

    loadProperties: function (lang) {

        jQuery.i18n.properties({
            name: 'lang', // 资源文件名称
            path: 'assets/i18n/', // 资源文件所在目录路径
            mode: 'map', // 模式：变量或 Map
            language: lang, // 对应的语言
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a').text($.i18n.prop('navbar_stake'));
                $('h3:eq(0)').text($.i18n.prop('stake_share_title'));
                $('.showShareDetail').text($.i18n.prop('stake_share_button'));
                $('.totalProfit div:eq(0)').text($.i18n.prop('stake_share_profit'));
                $('.totalShareNum div:eq(0)').text($.i18n.prop('stake_share_total'));
                $('.leftShareNum div:eq(0)').text($.i18n.prop('stake_share_left'));
                $('.hasShareNum div:eq(0)').text($.i18n.prop('stake_share_voted'));
                $('.missShareNum div:eq(0)').text($.i18n.prop('stake_share_missed'));
                $('.expireShareNum div:eq(0)').text($.i18n.prop('stake_share_expired'));
                $('h3:eq(1)').text($.i18n.prop('stake_pool_title'));
                $('.register').text($.i18n.prop('stake_pool_register'));
                $('thead tr td:eq(0)').text($.i18n.prop('stake_pool_id'));
                $('thead tr td:eq(1)').text($.i18n.prop('stake_pool_owner'));
                $('thead tr td:eq(2)').text($.i18n.prop('stake_pool_launched'));
                $('thead tr td:eq(3)').text($.i18n.prop('stake_pool_voted'));
                $('thead tr td:eq(4)').text($.i18n.prop('stake_pool_missed'));
                $('thead tr td:eq(5)').text($.i18n.prop('stake_pool_fee'));
                $('thead tr td:eq(6)').text($.i18n.prop('stake_pool_shareNum'));
                $('thead tr td:eq(7)').text($.i18n.prop('stake_pool_lastpaytime'));
                $('thead tr td:eq(8)').text($.i18n.prop('stake_pool_operation'));
                $('.buyShare').text($.i18n.prop('stake_pool_buyShare'));
            }
        });
    },

    getAccountlist: function () {

        var that = this;
        var biz = {}
        Common.post("account/list", biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var dataArray = res.biz;

                    var totalProfit = new BigNumber(0);
                    var expireShareNum = new BigNumber(0);
                    var leftShareNum = new BigNumber(0);
                    var missShareNum = new BigNumber(0);
                    var hasShareNum = new BigNumber(0);
                    var totalShareNum = new BigNumber(0);
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        Common.post('share/my', data.MainPKr, {}, function (res2) {
                            if (res2.base.code === 'SUCCESS') {
                                if (res2.biz.length > 0) {
                                    var data = res2.biz[0];

                                    var _totalProfit = new BigNumber(data.profit, 16).dividedBy(Common.baseDecimal);
                                    var _expireShareNum = new BigNumber(data.expired, 16);
                                    var _leftShareNum = new BigNumber(data.remaining, 16);
                                    var _missShareNum = new BigNumber(data.missed, 16);
                                    var _hasShareNum = new BigNumber(data.total, 16).minus(new BigNumber(data.missed)).minus(new BigNumber(data.remaining));
                                    var _totalShareNum = new BigNumber(data.total, 16);

                                    totalProfit = totalProfit.plus(_totalProfit);
                                    expireShareNum = expireShareNum.plus(_expireShareNum);
                                    leftShareNum = leftShareNum.plus(_leftShareNum);
                                    missShareNum = missShareNum.plus(_missShareNum);
                                    hasShareNum = hasShareNum.plus(_hasShareNum);
                                    totalShareNum = totalShareNum.plus(_totalShareNum);

                                    $('.totalProfit span:eq(1)').text(totalProfit.toFixed(6) + ' SERO')
                                    $('.expireShareNum span:eq(1)').text(expireShareNum.toString(10))
                                    $('.leftShareNum span:eq(1)').text(leftShareNum.toString(10))
                                    $('.missShareNum span:eq(1)').text(missShareNum.toString(10))
                                    $('.hasShareNum span:eq(1)').text(totalShareNum.minus(expireShareNum).minus(missShareNum).minus(leftShareNum).toString(10))
                                    $('.totalShareNum span:eq(1)').text(totalShareNum.toString(10))
                                }
                            }
                        });
                    }
                }
            }
        })
    },

    stakeList: function () {
        $('tbody').empty();
        Common.postAsync('stake', {}, {}, function (res) {
            if (res.base.code === 'SUCCESS') {

                var dataArray = res.biz;
                for (var data of dataArray) {
                    $('tbody').append(`
                    <tr>
                        <td class="text-break">${data.id}</td>
                        <td class="text-break">${data.own}</td>
                        <td>${data.closed ? "Closed" : "Opening"}</td>
                        <td>${new BigNumber(data.choicedNum, 16).minus(new BigNumber(data.missedNum, 16)).toString(10)}</td>
                        <td>${new BigNumber(data.missedNum, 16).toString(10)}</td>
                        <td>${new BigNumber(data.fee, 16).div(100).toString(10)}%</td>
                        <td>${new BigNumber(data.shareNum, 16).toString()}</td>
                        <td>${new BigNumber(data.lastPayTime, 16).toString(10)}</td>
                        <td><button class="btn btn-outline-primary btn-block small buyShare" attpoolid="${data.id}">${$.i18n.prop('stake_pool_buyShare')}</button></td>
                    </tr>
               `);
                }
            }
        });
    },


};

var StakeRegister = {

    init: function () {
        var that = this;

        that.getAccountlist();

        $('.close').bind('click', function () {
            $('.modal').hide();
            $('.modal-footer button:eq(1)').unbind('click');
            $("#sub1").attr('disabled', false);
        });


        $('.toast').toast({animation: true, autohide: true, delay: 2000})
    },

    loadProperties: function (lang) {

        jQuery.i18n.properties({
            name: 'lang', // 资源文件名称
            path: 'assets/i18n/', // 资源文件所在目录路径
            mode: 'map', // 模式：变量或 Map
            language: lang, // 对应的语言
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a').text($.i18n.prop('navbar_stake'));
                $('.breadcrumb li:eq(0) a').text($.i18n.prop('stake_pool_title'));
                $('.breadcrumb li:eq(1)').text($.i18n.prop('stake_pool_register'));
                $('h4').text($.i18n.prop('stake_pool_register'));
                $('.form-group:eq(0) label').text($.i18n.prop('stake_register_from'));
                $('.form-group:eq(1) label').text($.i18n.prop('stake_register_address'));
                $('.form-group:eq(1) .invalid-feedback').text($.i18n.prop('stake_register_address_tips'));
                $('.form-group:eq(2) label').text($.i18n.prop('stake_register_fee'));
                $('.form-group:eq(2) .invalid-feedback').text($.i18n.prop('stake_register_fee_tips'));
                $('.form-group:eq(3) label').text($.i18n.prop('stake_register_amount'));
                $('.form-group:eq(3) small').text($.i18n.prop('stake_register_amount_desc'));
                $('.form-group:eq(2) div div:eq(2)').text($.i18n.prop('stake_register_fee_desc'));
                $('.modal-title').text($.i18n.prop('stake_register_confirm_title'));
                $('.modal-body ul li:eq(0) div div:eq(0)').text($.i18n.prop('stake_register_from'));
                $('.modal-body ul li:eq(1) div div:eq(0)').text($.i18n.prop('stake_register_address'));
                $('.modal-body ul li:eq(2) div div:eq(0)').text($.i18n.prop('stake_register_fee'));
                $('.modal-body ul li:eq(3) div div:eq(0)').text($.i18n.prop('stake_register_amount'));
                $('.modal-body ul li:eq(4) div div:eq(0)').text($.i18n.prop('stake_register_password'));
                $('#password').attr('placeholder', $.i18n.prop('stake_register_password_place'));
                $('#sub1').text($.i18n.prop('stake_register_next'));
                $('#address').attr('placeholder', $.i18n.prop('stake_register_address_tips'));
                $('#feeRate').attr('placeholder', $.i18n.prop('stake_register_fee_tips'));

                $('.modal-footer button:eq(0)').text($.i18n.prop('stake_register_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('stake_register_confirm'));
            }
        });
    },

    getAccountlist: function () {

        var biz = {}
        Common.postAsync("account/list", biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var dataArray = res.biz;
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        var balance = new BigNumber(0).toFixed(6);
                        if (data.Balance) {
                            var balanceObj = data.Balance;
                            for (var currency of Object.keys(balanceObj)) {
                                if (currency === 'SERO') {
                                    balance = new BigNumber(balanceObj[currency]).dividedBy(Common.baseDecimal).toFixed(6);
                                    $('.address').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8) }   ${ balance + ' ' + currency}</option>`);
                                }
                            }
                        } else {
                            $('.address').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8) }   ${ '0.000 SERO'}</option>`);
                        }
                    }
                }
            }
        })
    },


    confirm: function () {

        var from = $(".address").val();
        var vote = $("#address").val();
        var feeRate = $("#feeRate").val();


        $('.modal-footer button:eq(0)').bind('click', function () {
            $('#sub1').attr('disabled', false);
            $('.modal-footer button:eq(1)').unbind('click');
        });

        $('ul:eq(1) li:eq(0) div div:eq(1)').text(from);
        $('ul:eq(1) li:eq(1) div div:eq(1)').text(vote);
        $('ul:eq(1) li:eq(2) div div:eq(1)').text(feeRate + '%');
        $('ul:eq(1) li:eq(3) div div:eq(1)').text("200,000 SERO");
        $('.modal').modal('show');

        $('.modal-footer button:eq(1)').bind('click', function () {
            var password = $("#password").val();
            if(password === ''){
                $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('stake_register_password_place'));
                $('.toast').toast('show');
            }else{
                $('.modal-footer button:eq(1)').attr('disabled', true).text($.i18n.prop('send_tx_sending'));
                var biz = {
                    From: from,
                    Vote: vote,
                    FeeRate: new BigNumber(feeRate).multipliedBy(100).toString(10),
                    Password: password,
                }
                Common.postAsync('stake/register', biz, {}, function (res) {
                    if (res.base.code === 'SUCCESS') {
                        $("#password").val('');
                        $('.toast-body').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                        $('.toast').toast('show');

                        $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                        $('#sub1').attr('disabled', false);
                    } else {
                        $('.toast-body').removeClass('alert-success').addClass('alert-danger').text(res.base.desc);
                        $('.toast').toast('show');
                        $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                    }
                })
            }

        });
    }


};

var StakeBuyer = {

    poolId: '',

    init: function () {
        var that = this;

        that.getAccountlist();

        that.poolId = GetQueryString('id');

        $('#poolId').val(that.poolId);

        $('.close').bind('click', function () {
            $('.modal').hide();
            $('.modal-footer button:eq(1)').unbind('click');
            $("#sub1").attr('disabled', false);
        });

        $('.register').bind('click', function () {
            window.location.href = 'stake-register.html';
        });
        $('.toast').toast({animation: true, autohide: true, delay: 2000})
    },

    loadProperties: function (lang) {

        jQuery.i18n.properties({
            name: 'lang', // 资源文件名称
            path: 'assets/i18n/', // 资源文件所在目录路径
            mode: 'map', // 模式：变量或 Map
            language: lang, // 对应的语言
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a').text($.i18n.prop('navbar_stake'));

                $('.breadcrumb li:eq(0) a').text($.i18n.prop('stake_pool_title'));
                $('.breadcrumb li:eq(1)').text($.i18n.prop('stake_pool_buyShare'));
                $('.form-group:eq(0) label').text($.i18n.prop('share_buy_poolId'));
                $('.form-group:eq(1) label').text($.i18n.prop('share_buy_from'));
                $('.form-group:eq(2) label').text($.i18n.prop('share_buy_address'));
                $('.form-group:eq(2) .invalid-feedback').text($.i18n.prop('share_buy_address_tips'));
                $('.form-group:eq(3) label').text($.i18n.prop('share_buy_amount'));
                $('.form-group:eq(3) .invalid-feedback').text($.i18n.prop('share_buy_amount_tips'));
                $('#address').attr('placeholder', $.i18n.prop('share_buy_address_tips'));
                $('#amount').attr('placeholder', $.i18n.prop('share_buy_amount_tips'));
                $('h4').text($.i18n.prop('stake_pool_buyShare'));
                $('.modal-title').text($.i18n.prop('stake_register_confirm_title'));
                $('.modal-body ul li:eq(0) div div:eq(0)').text($.i18n.prop('share_buy_poolId'));
                $('.modal-body ul li:eq(1) div div:eq(0)').text($.i18n.prop('share_buy_from'));
                $('.modal-body ul li:eq(2) div div:eq(0)').text($.i18n.prop('share_buy_address'));
                $('.modal-body ul li:eq(3) div div:eq(0)').text($.i18n.prop('share_buy_amount'));
                $('.modal-body ul li:eq(4) div div:eq(0)').text($.i18n.prop('stake_register_password'));
                $('#password').attr('placeholder', $.i18n.prop('stake_register_password_place'));
                $('#sub1').text($.i18n.prop('stake_register_next'));

                $('.modal-footer button:eq(0)').text($.i18n.prop('stake_register_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('stake_register_confirm'));
            }
        });
    },

    getAccountlist: function () {

        var biz = {}
        Common.postAsync("account/list", biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var dataArray = res.biz;
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        var balance = new BigNumber(0).toFixed(6);
                        if (data.Balance) {
                            var balanceObj = data.Balance;
                            for (var currency of Object.keys(balanceObj)) {
                                if (currency === 'SERO') {
                                    balance = new BigNumber(balanceObj[currency]).dividedBy(Common.baseDecimal).toFixed(6);
                                    $('.address').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8) }   ${ balance + ' ' + currency}</option>`);
                                }
                            }
                        } else {
                            $('.address').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8) }   ${ '0.000 SERO'}</option>`);
                        }
                    }
                }
            }
        })
    },


    confirm: function () {

        var from = $(".address").val();
        var vote = $("#address").val();
        var amount = $("#amount").val();
        var poolId = $("#poolId").val();


        $('.modal-footer button:eq(0)').bind('click', function () {
            $('#sub1').attr('disabled', false);
            $('.modal-footer button:eq(1)').unbind('click');
        });

        $('ul:eq(1) li:eq(0) div div:eq(1)').text(poolId);
        $('ul:eq(1) li:eq(1) div div:eq(1)').text(from);
        $('ul:eq(1) li:eq(2) div div:eq(1)').text(vote);
        $('ul:eq(1) li:eq(3) div div:eq(1)').text(amount);
        $('.modal').modal('show');

        $('.modal-footer button:eq(1)').bind('click', function () {
            var password = $("#password").val();
            if(password === ''){
                $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('send_tx_success'));
                $('.toast').toast('show');
            }else{
                $('.modal-footer button:eq(1)').attr('disabled', true);
                var biz = {
                    From: from,
                    Vote: vote,
                    Amount: new BigNumber(amount).multipliedBy(Common.baseDecimal).toString(10),
                    Pool: poolId,
                    Password: password,
                    GasPrice: new BigNumber(1000000000).toString(10),
                }
                Common.postAsync('stake/buyShare', biz, {}, function (res) {
                    if (res.base.code === 'SUCCESS') {
                        $('.toast-body').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                        $('.toast').toast('show');
                        $('.modal-footer button:eq(1)').attr('disabled', false);
                        $('#sub1').attr('disabled', false);
                        setTimeout(function () {
                            window.location.href = 'stake.html';
                        }, 2000);
                    } else {
                        $('.toast-body').removeClass('alert-success').addClass('alert-danger').text(res.base.desc);
                        $('.toast').toast('show');
                        $('.modal-footer button:eq(1)').attr('disabled', false);
                    }
                })
            }
        });
    }


};


var StakeDetail = {


    init: function () {
        var that = this;

        that.getAccountlist();

    },

    loadProperties: function (lang) {

        jQuery.i18n.properties({
            name: 'lang', // 资源文件名称
            path: 'assets/i18n/', // 资源文件所在目录路径
            mode: 'map', // 模式：变量或 Map
            language: lang, // 对应的语言
            cache: false,
            encoding: 'UTF-8',
            callback: function () {
                $('.navbar-nav li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a').text($.i18n.prop('navbar_stake'));
                $('.shareDetail').text($.i18n.prop('share_detail_title'));
                $('thead tr td:eq(0) strong').text($.i18n.prop('share_detail_shareId'));
                $('thead tr td:eq(1) strong').text($.i18n.prop('share_detail_poolId'));
                $('thead tr td:eq(2) strong').text($.i18n.prop('share_detail_address'));
                $('thead tr td:eq(3) strong').text($.i18n.prop('share_detail_price'));
                $('thead tr td:eq(4) strong').text($.i18n.prop('share_detail_fee'));
                $('thead tr td:eq(5) strong').text($.i18n.prop('share_detail_profit'));
                $('thead tr td:eq(6) strong').text($.i18n.prop('share_detail_remaining'));
                $('thead tr td:eq(7) strong').text($.i18n.prop('share_detail_voted'));
                $('thead tr td:eq(8) strong').text($.i18n.prop('share_detail_expired'));
                $('thead tr td:eq(9) strong').text($.i18n.prop('share_detail_missed'));
                $('thead tr td:eq(10) strong').text($.i18n.prop('share_detail_total'));

                $('.breadcrumb li:eq(0) a').text($.i18n.prop('stake_pool_title'));
                $('.breadcrumb li:eq(1)').text($.i18n.prop('share_detail_title'));


            }
        });
    },

    getAccountlist: function () {

        var biz = {}
        Common.post("account/list", biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var dataArray = res.biz;

                    $('tbody').empty();
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        Common.post('share/my', data.MainPKr, {}, function (res) {
                            if (res.base.code === 'SUCCESS') {
                                if (res.biz.length > 0) {
                                    var data = res.biz[0];
                                    var shareIds = data.shareIds;
                                    for (let shareId of shareIds) {
                                        Common.post('stake/getShare', shareId, {}, function (res) {
                                            var share = res.biz;
                                            if (res.base.code === 'SUCCESS') {
                                                $('tbody').append(`
                                                <tr>
                                                <td>${share.id.substring(0, 5) + " ... " + share.id.substring(share.id.length - 5)}</td>
                                                <td>${share.pool.substring(0, 5) + " ... " + share.pool.substring(share.pool.length - 5)}</td>
                                                <td>${share.addr.substring(0, 5) + " ... " + share.addr.substring(share.addr.length - 5)}</td>
                                                <td>${new BigNumber(share.price, 16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                                <td>${(parseFloat(new BigNumber(share.fee,16).toString(10)) / 10000).toFixed(2)}</td>
                                                <td>${new BigNumber(share.profit, 16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                                <td>${new BigNumber(share.remaining?share.remaining:"0x0", 16).toString(10)}</td>
                                                <td>${new BigNumber(share.total, 16).minus(new BigNumber(share.remaining?share.remaining:"0x0", 16)).minus(new BigNumber(share.missed?share.missed:"0x0", 16)).minus(new BigNumber(share.expired?share.expired:"0x0", 16))}</td>
                                                <td>${share.expired ? new BigNumber(share.expired, 16).toString(10) : "0"}</td>
                                                <td>${new BigNumber(share.missed?share.missed:"0x0", 16).toString(10)}</td>
                                                <td>${new BigNumber(share.total,16).toString(10)}</td>
                                                </tr>
                                            `);
                                            }
                                        });
                                    }
                                }
                            }
                        });
                    }
                }
            }
        })
    },
}

function GetQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]);
    return null;
}