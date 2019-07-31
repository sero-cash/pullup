var Account = {


    init: function () {


        $('.close').bind('click', function () {
            $('.modal').hide();
        });

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
                $('p:eq(0)').text($.i18n.prop('account_new_tips'));
                $('h4:eq(0)').text($.i18n.prop('account_new'));
                $('input:eq(0)').attr('placeholder',$.i18n.prop('account_new_password_place')).next('div').text($.i18n.prop('account_new_password_place'));;
                $('input:eq(1)').attr('placeholder',$.i18n.prop('account_new_password_confirm')).next('div').text($.i18n.prop('account_new_password_place'));
                $('button:eq(0)').text($.i18n.prop('account_new_back'));
                $('button:eq(1)').text($.i18n.prop('account_new_next'));

                // $('a:eq(0) span').text($.i18n.prop('account_import_keystore'));
                $('a:eq(0) span').text($.i18n.prop('account_import_mnemnic'));

                $('.modal-title').text($.i18n.prop('account_new_modal_title'));
                $('.modal-body span:eq(0)').text($.i18n.prop('account_new_modal_address'));
                $('.modal-body span:eq(1)').text($.i18n.prop('account_new_modal_mnemonic'));
                $('.modal-footer button:eq(0)').text($.i18n.prop('account_new_modal_close'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('account_new_modal_access'));
            }
        });
    },

    newAccount: function () {

        var pwd = $("#pwd").val();
        var confirmPwd = $("#confirmPwd").val();

        if (pwd !== confirmPwd) {
            $('.toast-body').text($.i18n.prop('account_new_modal_pwdtips'));
            $('.toast').toast('show');
            $('#sub1').attr('disabled', false);
        } else {
            var biz = {
                passphrase: pwd,
            }
            Common.post("account/create", biz, {}, function (res) {
                if (res.base.code === "SUCCESS") {
                    $('.modal-title').text($.i18n.prop('account_new_modal_title'));
                    $('#myModal').modal({backdrop: 'static', keyboard: false})
                    $('.modal-footer button:eq(0)').bind('click', function () {
                        $("#sub1").attr('disabled', false);
                        $('.modal-body p:eq(1)').text('');
                        $('#myModal').modal('hide');
                    });
                    $('.modal-footer button:eq(1)').bind('click', function () {
                        window.location.href = "index.html";
                    });
                    $('.modal-body p:eq(0)').text(res.biz.address.substring(0,20) + " ... " + res.biz.address.substring(res.biz.address.length-20));
                    $('.modal-body p:eq(1)').text(res.biz.mnemonic);
                    $('#myModal').modal({backdrop: 'static', keyboard: false});
                    $("#sub1").attr('disabled', false);
                } else {
                    $("#sub1").text("NEXT").attr('disabled', false);
                    alert(res.base.desc);
                }

            })
        }
    },
};

var Detail = {

    currencyDecimal:{},

    address: '',

    txList: [],

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
                $('.breadcrumb li:eq(0) a').text($.i18n.prop('navbar_home'));
                $('.breadcrumb li:eq(1)').text($.i18n.prop('account_detail'));
                $('h3:eq(0)').text($.i18n.prop('account_detail'));
                $('h3:eq(1)').text($.i18n.prop('account_balance'));
                $('.c_a').text($.i18n.prop('account_collectionAddress'));
                $('.mainPkr').text($.i18n.prop('account_mainPkr'));
                $('.main-pkr-warning').text($.i18n.prop('account_mainPkrDesc'));
                $('.pkr-warning').text($.i18n.prop('account_pkrDesc'));

                $('.tx-his').text($.i18n.prop('account_tx'));
                $('table tr td:eq(1) strong').text($.i18n.prop('account_tx_hash'));
                $('table tr td:eq(2) strong').text($.i18n.prop('account_tx_block'));
                $('table tr td:eq(3) strong').text($.i18n.prop('account_tx_currecy'));
                $('table tr td:eq(4) strong').text($.i18n.prop('account_tx_state'));
                $('table tr td:eq(5) strong').text($.i18n.prop('account_tx_amount'));
                $('table tr td:eq(6) strong').text($.i18n.prop('account_tx_fee'));
                $('.network').text($.i18n.prop('home_account_network'));
                $('.addAcount').text($.i18n.prop('home_account_add'));
                $('.exportPhrase').text($.i18n.prop('account_export_mnemnic'));
                $('.modal-footer button:eq(0)').text($.i18n.prop('send_tx_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('send_tx_confirm'));
            }
        });
    },

    init: function () {
        var that = this;
        that.getAccountDetail();
        that.getTxList();

        $('.toast').toast({animation: true, autohide: true, delay: 1000})
        var clipboard1 = new ClipboardJS('.fa-copy');
        clipboard1.on('success', function (e) {
            $('#toast1 div:eq(0)').text('Copy successfully!');
            $('#toast1').toast('show')
        });


        that.bindExport();

        setInterval(function () {
            that.getAccountDetail();
            that.getTxList();
        }, 20000);
    },

    bindExport:function () {
        var that = this;
        $('.backup').bind('click', function () {
            $('.modal-title').empty().text($.i18n.prop('account_export_mnemnic'));
            $('.modal-body div:eq(1)').empty().text('');
            $('#qrcode').empty().append(`
                <input type="password" id="password" class="form-control" placeholder="${$.i18n.prop('send_tx_pwdtips')}" maxlength="50" aria-describedby="inputGroup-sizing-sm">
            `);
            $('.modal').modal('show');
            $('.modal-footer button:eq(1)').unbind('click').bind('click',function () {
                var pasword = $('#password').val();
                if (pasword){
                    $('.modal-footer button:eq(1)').attr('disabled',true);
                    $('.modal-body div:eq(1)').empty();
                    var biz ={
                        passphrase:pasword,
                        address:that.address,
                    }
                    Common.postAsync('account/export/mnemonic', biz, {}, function (res) {
                        if (res.base.code === 'SUCCESS') {
                            $('.modal-title').empty().text($.i18n.prop('account_new_modal_title'));
                            $('#qrcode').empty().append(`<span class="text-dark">${res.biz}</span>`);
                            $('.modal-footer button:eq(1)').unbind('click').bind('click',function () {
                                $('.modal').modal('hide');
                                $('.modal-body p:eq(1)').text('');
                            }).attr('disabled',false);
                        }else{
                            $('.modal-body div:eq(1)').empty().append(`<strong class="text-danger">${res.base.desc}</strong>`);
                            $('.modal-footer button:eq(1)').attr('disabled',false);
                        }

                    });
                }else{
                    $('.modal-body div:eq(1)').empty().append(`<span class="text-danger">${$.i18n.prop('send_tx_pwdtips')}</span>`);
                }

            });
        });
    },

    getAccountDetail: function () {

        var that = this;
        var pk = GetQueryString("pk");
        var biz = {
            PK: pk,
        }
        $('.currency').empty();
        Common.post("account/detail", biz, {}, function (res) {

            if (res.base.code === "SUCCESS") {

                that.address = res.biz.PK;
                var pkr = res.biz.PkrBase58;

                if (res.biz.Name){
                    $('.a_span').text(res.biz.Name);
                }

                $('.mainCopy').attr('data-clipboard-text', res.biz.MainPKr);
                $('.secondeCopy').attr('data-clipboard-text', pkr);
                $('.address').text(pkr);
                $('.main-address').text(res.biz.MainPKr);


                $('.mainqrcode').bind('click', function () {
                    $('.modal-title').empty().text("Qrcode");
                    $('.modal-body div:eq(1)').empty().text(res.biz.MainPKr);
                    $('#qrcode').empty().qrcode({
                        render: "canvas",
                        width: 200,
                        height: 200,
                        text: res.biz.MainPKr
                    });
                    $('.modal-footer button:eq(1)').bind('click',function () {
                        $('#myModal').modal('hide');
                        $('.modal-footer button:eq(1)').unbind('click');
                    })
                    $('#myModal').modal({backdrop: 'static', keyboard: false});
                });

                $('.secondqrcode').bind('click', function () {
                    $('.modal-title').empty().text("Qrcode");
                    $('.modal-body div:eq(1)').empty().text(pkr);
                    $('#qrcode').empty().qrcode({
                        render: "canvas",
                        width: 200,
                        height: 200,
                        text: pkr
                    });
                    $('.modal-footer button:eq(1)').bind('click',function () {
                        $('#myModal').modal('hide');
                        $('.modal-footer button:eq(1)').unbind('click');
                    })
                    $('#myModal').modal('show');
                });

                $('.pk').text(res.biz.PK.substring(0, 8) + " ... " + res.biz.PK.substring(res.biz.PK.length - 8));

                var balanceObj = res.biz.Balance;

                var strMap = new Map();
                for (var k of Object.keys(balanceObj)) {
                    strMap.set(k, balanceObj[k]);
                    if (k !== 'SERO'){
                        var biz={
                            Currency:k,
                        }
                        Common.post('decimal',biz,{},function (res) {
                            var cDecimal = new BigNumber(10).pow(new BigNumber(res.biz));
                            var amount = new BigNumber(balanceObj[k]).dividedBy(cDecimal);
                            if(res.biz>6){
                                amount = amount.toFixed(6);
                            }else{
                                amount = amount.toFixed(res.biz);
                            }
                            $('.currency').append(`
                                <div class="col-md-3 col-xl-3 mb-4">
                                    <div class="card shadow border-left-success py-2">
                                        <div class="card-body">
                                            <div class="row align-items-center no-gutters">
                                                <div class="col mr-2">
                                                    <div class="text-uppercase text-success font-weight-bold text-xs mb-1"><span>${k}</span></div>
                                                    <div class="text-dark font-weight-bold h5 mb-0"><span>${amount}</span></div>
                                                </div>
                                                <div class="col-auto"><i class="fas fa-dollar-sign fa-2x text-gray-300"></i></div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            `);
                        });
                    }else{
                        $('.currency').append(`
                        <div class="col-md-3 col-xl-3 mb-4">
                            <div class="card shadow border-left-success py-2">
                                <div class="card-body">
                                    <div class="row align-items-center no-gutters">
                                        <div class="col mr-2">
                                            <div class="text-uppercase text-success font-weight-bold text-xs mb-1"><span>${k}</span></div>
                                            <div class="text-dark font-weight-bold h5 mb-0"><span>${new BigNumber(balanceObj[k]).dividedBy(Common.baseDecimal).toFixed(6)}</span></div>
                                        </div>
                                        <div class="col-auto"><i class="fas fa-dollar-sign fa-2x text-gray-300"></i></div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    `);
                    }
                }

                if (strMap.size === 0) {

                    $('.currency').append(`
                        <div class="col-md-3 col-xl-3 mb-4">
                            <div class="card shadow border-left-success py-2">
                                <div class="card-body">
                                    <div class="row align-items-center no-gutters">
                                        <div class="col mr-2">
                                            <div class="text-uppercase text-success font-weight-bold text-xs mb-1"><span>SERO</span></div>
                                            <div class="text-dark font-weight-bold h5 mb-0"><span>0.000000</span></div>
                                        </div>
                                        <div class="col-auto"><i class="fas fa-dollar-sign fa-2x text-gray-300"></i></div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    `);
                }

            }
        });
    },

    txPageNo: 1,

    txPageSize: 10,

    txCount: 0,


    getTxList: function (callback) {
        var that = this;

        var pk = GetQueryString("pk");
        var biz = {
            PK: pk,
        };
        var page = {
            page_no: that.txPageNo,
            page_size: that.txPageSize,
        }
        $('tbody').empty();

        Common.post("tx/list", biz, page, function (res) {
            if (res.base.code === "SUCCESS") {
                if (res.biz) {
                    var data = res.biz;
                    for (var i = 0; i < data.length; i++) {

                        var tx = data[i];
                        var amount = new BigNumber(tx.Amount);


                        if(tx.Currency !== 'SERO'){
                            var biz={
                                Currency:tx.Currency,
                            }
                            Common.post('decimal',biz,{},function (res) {
                                var cDecimal = new BigNumber(10).pow(new BigNumber(res.biz));
                                if(res.biz>6){
                                    amount = amount.dividedBy(cDecimal).toFixed(6);
                                }else{
                                    amount = amount.dividedBy(cDecimal).toFixed(res.biz);
                                }
                                $('tbody').append(
                                    `
                                    <tr>
                                        <td>${i + 1}</td>
                                        <td class="text-info"><a target="_blank" href="https://explorer.web.sero.cash/txsInfo.html?hash=${tx.Hash}">${tx.Hash.substring(0, 5) + " ... " + tx.Hash.substring(tx.Hash.length - 5)}</a></td>
                                        <td>${tx.Block==999999999999999?0:tx.Block}</td>
                                        <!--<td title="${tx.To}">${tx.To.substring(0, 5) + " ... " + tx.To.substring(tx.To.length - 5)}</td>-->
                                        <td>${tx.Currency}</td>
                                        <td><span class="text-success">${tx.Block!=999999999999999 ? 'Completed' : 'Pending'}</span></td>
                                        <td>${amount}</td>
                                        <td>${new BigNumber(tx.Fee).dividedBy(Common.baseDecimal).toFixed(8)}</td>
                                    </tr>
                                    `
                                            );
                            });
                        }else{
                            if( amount.comparedTo(new BigNumber(0))<0 && tx.Block !=999999999999999){
                                var fee = new BigNumber(tx.Fee)
                                if(amount.plus(fee) < 0 ){
                                    amount = amount.plus(fee)
                                }else{
                                    tx.Fee = 0;
                                }
                            }
                            amount = amount.dividedBy(Common.baseDecimal).toFixed(6);
                            $('tbody').append(
                                `
                            <tr>
                                <td>${i + 1}</td>
                                <td class="text-info"><a target="_blank" href="https://explorer.web.sero.cash/txsInfo.html?hash=${tx.Hash}">${tx.Hash.substring(0, 5) + " ... " + tx.Hash.substring(tx.Hash.length - 5)}</a></td>
                                <td>${tx.Block==999999999999999?0:tx.Block}</td>
                                <!--<td title="${tx.To}">${tx.To.substring(0, 5) + " ... " + tx.To.substring(tx.To.length - 5)}</td>-->
                                <td>${tx.Currency}</td>
                                <td><span class="text-success">${tx.Block !=999999999999999 ? 'Completed' : 'Pending'}</span></td>
                                <td>${amount}</td>
                                <td>${new BigNumber(tx.Fee).dividedBy(Common.baseDecimal).toFixed(8)}</td>
                            </tr>
                            `
                                );
                        }

                    };
                }
                if (callback) {
                    callback(res.page)
                }
            }
        });
    }
}


var Keystore = {

    file: '',

    init: function () {
        var that = this;

        $('.close').bind('click', function () {
            $('#myModal').modal('hide');
        });

        $('.modal-footer button:eq(1)').bind('click', function () {
            window.location.href = 'index.html';
        });


        $("#i-file").bind("change", function () {
            that.file = this.files[0];
        });
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
                $('h3:eq(0)').text($.i18n.prop('home_account'));
                $('.total').text($.i18n.prop('home_account_total'));
                $('.blockHeight').text($.i18n.prop('home_account_height'));
                $('.network').text($.i18n.prop('home_account_network'));
                $('.addAcount').text($.i18n.prop('home_account_add'));
            }
        });
    },

    import: function () {
        var that = this;
        var password = $('#password').val();
        var formData = new FormData();
        formData.append("passphrase", password);
        formData.append("uploadFile", that.file);

        $.ajax({
            url: Common.host + "/account/import/keystore",
            dataType: 'json',
            type: 'POST',
            async: false,
            data: formData,
            processData: false,
            contentType: false,
            success: function (data) {
                if (data.responseText === 'INVALID_FILE_TYPE') {
                    $('.modal-title').text("Warning");
                    $('.modal-body').text("Password given is incorrect!");
                } else if (data.responseText === 'SUCCESS') {
                    $('.modal-title').text("Successful");
                    $('.modal-body').text("Successfully imported!");
                } else {
                    $('.modal-title').text("Error");
                    $('.modal-body').text("Import failed,Incorrect file type");
                }
                $('#myModal').modal({backdrop: 'static', keyboard: false});
            },
            error: function (data) {
                if (data.responseText === 'INVALID_FILE_TYPE') {
                    $('.modal-title').text("Warning");
                    $('.modal-body').text("Password given is incorrect!");
                } else if (data.responseText === 'SUCCESS') {
                    $('.modal-title').text("Successful");
                    $('.modal-body').text("Successfully imported!");
                } else {
                    $('.modal-title').text("Error");
                    $('.modal-body').text("Import failed,Incorrect file type");
                }
                $('#myModal').modal('show');
            }
        });

        $("#sub1").attr('disabled', false);
    },


};

var Mnemnic = {

    init: function () {
        $('.close').bind('click', function () {
            $('.modal').modal('hide');
        });

        $('.modal-footer button:eq(1)').bind('click', function () {
            window.location.href = 'index.html';
        });

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
                //$('p:eq(0)').text($.i18n.prop('navbar_home'));

            }
        });
    },

    import: function () {
        var mnemonic = $('#mnemonic').val();
        var password = $('#password').val();

        var biz = {
            mnemonic: mnemonic,
            passphrase: password,
        }

        Common.post('account/import/mnemonic', biz, {}, function (res) {

            if (res.base.code === 'SUCCESS') {
                var address = res.biz.address
                $('.modal-title').text("Import Successful");
                $('.modal-body p:eq(0)').text(address.substring(0,20) + " ... " + address.substring(address.length-20));
            } else {
                $('.modal-title').text("ERROR");
                $('.modal-body p:eq(0)').text(res.base.desc);
            }
            $('#myModal').modal({backdrop: 'static', keyboard: false});
            $("#sub1").attr('disabled', false);
        });

    }

}



function GetQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if (r != null) return unescape(r[2]);
    return null;
}
