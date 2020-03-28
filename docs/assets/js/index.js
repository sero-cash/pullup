var Index = {

    network: {},

    init: function () {
        var that = this;
        that.getAccountList();
        that.getBlockHeight();
        setInterval(function () {
            that.getAccountList();
        }, 10000);

        setInterval(function () {
            that.getBlockHeight();
        }, 10000);

        $('.select-net').bind('click', function () {
            that.selectNetwork();
        });

        $('.backup').bind('click',function () {
            that.backup();
        });

        Common.post('network/change', "", {}, function (res) {
            if (res.base.code === 'SUCCESS') {

                if(res.biz === " http://140.143.83.98:8545" ){
                    Common.post('network/change', "http://140.143.83.98:8545", {}, function (res) {
                        if (res.base.code === 'SUCCESS') {
                            $.cookie('networkUrl', res.biz);
                            $.cookie('seroRpcHost',res.biz);
                            $('.select-net span').text(res.biz);
                        }
                    });
                }else{
                    $.cookie('networkUrl', res.biz);
                    $.cookie('seroRpcHost',res.biz);
                    $('.select-net span').text(res.biz);
                }
            }
        });
    },

    backup:function () {
        Common.post('file/open',{},{},function (res) {
        });
    },

    selectNetwork: function () {

        var localNet = $.cookie('networkUrl');
        $('.list-group').empty();

        $.getJSON("node.json", "", function(data) {　
            $.each(data.host, function(i, net) {
                $('.list-group').append(`<li class="list-group-item ${localNet === net.rpc ? 'active' : ''} netcheck" network="${net.rpc}" netname="${net.name}">
                    <i class="fas fa-circle ${net.network === 'main' ? 'text-success' : 'text-warning'}"></i> ${net.rpc} ${net.name}
                </li>`);
            })

            $('.list-group').append(`<li class="list-group-item" network="personal" netname="Personal Network" >
                <div class="input-group mb-3">
                <div class="input-group-prepend">
                <span class="input-group-text" id="basic-addon3">http://</span>
                </div>
                <input type="text" class="form-control" id="basic-url" aria-describedby="basic-addon3" placeholder="127.0.0.1:8545">
                </div>
                <button class="btn btn-outline-primary btn-block addNetwork">Set Personal RPC</button>
            </li>`);

            $('.netcheck').bind('click', function () {
                var network = $(this).attr('network');
                // var netName = $(this).attr('netname');
                Common.post('network/change', network, {}, function (res) {
                    if (res.base.code === 'SUCCESS') {
                        $('.select-net span').text(network);
                        $('.netcheck').unbind('click');
                        $('#myModal').modal('hide');
                        // $.cookie('networkName', netName);
                        $.cookie('networkUrl', network);
                        $.cookie('seroRpcHost',res.biz);
                    }
                });
            });

            $('.addNetwork').bind('click', function () {
                var network = $('#basic-url').val();
                // var netName = $(this).attr('netname');
                if (network === '') {
                    $('.toast-body').removeClass('alert-success').addClass('alert-danger').text("Please Enter Network");
                    $('.toast').toast('show');
                } else {
                    if ('/' === (network.substring(network.length - 1))) {
                        network = network.substring(network.length - 1);
                    }
                    network = "http://" + network;
                    Common.post('network/change', network, {}, function (res) {
                        if (res.base.code === 'SUCCESS') {
                            // $.cookie('networkName', network);
                            $.cookie('networkUrl', network);
                            $.cookie('seroRpcHost',res.biz);
                            $('.select-net span').text(network);
                            $('.netcheck').unbind('click')
                            $('.addNetwork').unbind('click');
                            $('.toast-body').removeClass('alert-danger').addClass('alert-success').text("Set Network Success");
                            $('.toast').toast('show');
                            setTimeout(function () {
                                $('#myModal').modal('hide');
                            }, 1000);
                        }
                    });
                }
            });
        });

        $('#myModal').modal({backdrop: 'static', keyboard: false});

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
                $('.navbar-nav li:eq(3) a').text($.i18n.prop('navbar_dapps'));

                $('h3:eq(0)').text($.i18n.prop('home_account'));
                $('.total').text($.i18n.prop('home_account_total'));
                $('.blockHeight').text($.i18n.prop('home_account_height'));
                $('.network').text($.i18n.prop('home_account_network'));
                $('.addAcount').text($.i18n.prop('home_account_add'));
            }
        });
    },

    getAccountList: function () {
        var biz = {}
        $('.pkrs').empty();
        Common.post("account/list", biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {

                if (res.biz) {

                    var dataArray = res.biz;
                    var balance = new BigNumber(0);

                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        var _balance = new BigNumber(0);
                        if (data.Balance && data.Balance.SERO) {
                            _balance = new BigNumber(data.Balance.SERO);
                            _balance = _balance.dividedBy(Common.baseDecimal);
                            balance = balance.plus(_balance)
                        }
                        var acName = "Account"+(i + 1);
                        if (data.Name){
                            acName = data.Name;
                        }
                        $('.pkrs').append(`
                            
                            <div class="col-lg-12 mb-4">
                                <div class="card text-white bg-primary shadow">
                                    <div class="card-body">
                                         <a style="text-decoration: none;color: white;" href="account-detail.html?pk=${data.PK}">
                                            <p class="m-0">${acName}(<small>${data.PK.substring(0, 8) + " ... " + data.PK.substring(data.PK.length - 8, data.PK.length)}</small>)</p>
                                            <p class="text-white-50 small m-0 pkr">
                                            ${data.PkrBase58}&nbsp;&nbsp &nbsp;</p>
                                            <p class="text-right text-warning m-0"><strong>${_balance.toFixed(6)}</strong> SERO</p>
                                         </a>
                                    </div>
                                </div>
                            </div>
                           
                        `);
                    }

                    $('.dashboard span:eq(0)').text(balance.toFixed(6));
                }
            }
        })
    },

    getBlockHeight: function () {

        Common.postAsync('sero/getBlockNumber', {}, {}, function (res) {
            if (res.base.code === 'SUCCESS') {
                if (res.biz) {
                    var blockNumber = new BigNumber(res.biz, 16).toString(10);
                    $('.dashboard span:eq(1)').text(blockNumber);
                }
            }
        })

    },

};