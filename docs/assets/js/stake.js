var StakeHome = {

    account:{},

    accountMap:{},

    stakeName:{
        zh_CN:{
            "0x388b2c9ba68a96bf697602fef9219f64e4ff8aef49815d0aeb56afd2a1276942":"赛罗",
            "0x82abc9d07aa976761cede08e53de8d5057efd81fc0c443c59b593231e69b4291":"SERDAC",
            "0xbde17513156dbbd0730b7dde954ce5d66930d646ac60a2f118572f56960c9d59":"盖世",
            "0xfeb23ac54e8d93994689bd782140b5804cfeec9d51e5d5986b35d0d843d1c146":"币龙驴池",
            "0x98f53bdad932c3865eebb229d0f74c4d2ee40440cfc2d34bf2ddec0a836f6f8d":"牛池",
            "0xc8db791edb4d2063f625de473a5061f9323114cb9d6de6bdfc82bbbba82642f0":"盘古",
            "0xc248ba3e8f98ec6714a9c3b59c4422cbc473b90c0d4fb01e589f5b8ae20a24d7":"马努",
            "0x16759fd13a7143207b3ebb088711b242267303dcdad53562d45fb4cfaf5dbdac":"仙龙节点",
            "0xda06d65e49808f31dec7b44339d856ff47ad2040a503ccd28a43a681195b23e1":"Hotbit ",
            "0x4fb40805a34c590cc78ca3d5e4f938a64424db4d4326e87d314a82e1d676bd60":"第一POS",
            "0xcec0343b0b29eecb24ec54dafcb97adfedc2acc367348b851e71973aa4e54659":"菠菜",
            "0xf1df2afb326a544a928a229a94f5eb8433d39688b590acd41c73d08200480b86":"雪庄Rose",
            "0xbdb9555b61613f8b13fd16918c9a09e407c3e96afdf8fe5dc887317eb0253cd7":"蚂蚁",
            "0x98d84dc25b65cf32a8488f04e728396fa96a15db682d79cde213a2368abb84d8":"HyperPay",
            "0x86fea8e64ee24e515846642822ab7c0fabaf3f2879fb3064f267047d300b76f3":"9池节点",
            "0xf26f8e060b241ad3ff9c5821d85544b21419ef2730c5c64dc7bd89bfc3710885":"T网矿池",
            "0x6cd71031cdc64db40816eecb01f65b3b56cdbde843bcd0797424e6184bf88e07":"魔完",
            "0xe058df83ef89df930aca3fed0c949ab6de1ccb51d2a87fc056631f11e3cd2ccc":"分水关",
            "0x24a0a75882b3be4eb1b5ceee17c6f253d147cfbc43cb35e798cd5560f43d93ab":"QBTC节点",
            "0x5c7dbbe680d7687e23b0dd927c37e77f051d701d57ef118f2e7ef5ad9bf50c51":"超零猫猫",
            "0x36e40dc58bf2db57f1831a02c2a248b0fa378ba3a99908c17a63ed622174d496":"超零猴哥",
            "0x16a800063ef8345f8441de438d97c8a49bc3373778e45a2f2cd25699be8d2c2b":"抹茶节点",
            "0xfd703f582147e4725a9b3567635dde7e442182c1f28e315360e07692a968183b":"超级节点",
            "0xb688861cf4b0f74509a4a4a0bdcdbe80ec0b5647bfcc41973498c9b8152fce41":"OKEx Pool",
            "0xce85c13d0942c3118873facf99669d193711bdca77b3de36668fdf29429bd8e5":"追日节点",
            "0x92af05b6d6c27a73eb8489f100c859cf3167298915cbae7262345e17aee6aa2c":"Bilaxy 节点",
            "0x7fbbc2f11ac02255b48c8ab534d08a0e2fa134d483794ca9294ad296f8e71df7":"HCOIN",
            "0x3800ad263e7b866d5fb21a59f867e415e6d8515b9ae58c84a9c1cdb38e7ba729":"BiKi矿池",

        },
        en_US:{
            "0x388b2c9ba68a96bf697602fef9219f64e4ff8aef49815d0aeb56afd2a1276942":"SERO Node (sero.im)",
            "0x82abc9d07aa976761cede08e53de8d5057efd81fc0c443c59b593231e69b4291":"SERDAC Node",
            "0xbde17513156dbbd0730b7dde954ce5d66930d646ac60a2f118572f56960c9d59":"Galaxy Node",
            "0xfeb23ac54e8d93994689bd782140b5804cfeec9d51e5d5986b35d0d843d1c146":"Coin Dragon Pool Node",
            "0x98f53bdad932c3865eebb229d0f74c4d2ee40440cfc2d34bf2ddec0a836f6f8d":"Newbit Node",
            "0xc8db791edb4d2063f625de473a5061f9323114cb9d6de6bdfc82bbbba82642f0":"Pangu Node",
            "0xc248ba3e8f98ec6714a9c3b59c4422cbc473b90c0d4fb01e589f5b8ae20a24d7":"Manu Node",
            "0x16759fd13a7143207b3ebb088711b242267303dcdad53562d45fb4cfaf5dbdac":"Immortal Dragon Node",
            "0xda06d65e49808f31dec7b44339d856ff47ad2040a503ccd28a43a681195b23e1":"Hotbit Node ",
            "0x4fb40805a34c590cc78ca3d5e4f938a64424db4d4326e87d314a82e1d676bd60":"The First PoS Node",
            "0xcec0343b0b29eecb24ec54dafcb97adfedc2acc367348b851e71973aa4e54659":"Spinach Node",
            "0xf1df2afb326a544a928a229a94f5eb8433d39688b590acd41c73d08200480b86":"Snow Rose Node",
            "0xbdb9555b61613f8b13fd16918c9a09e407c3e96afdf8fe5dc887317eb0253cd7":"Ant Node",
            "0x98d84dc25b65cf32a8488f04e728396fa96a15db682d79cde213a2368abb84d8":"HyperPay Node",
            "0x86fea8e64ee24e515846642822ab7c0fabaf3f2879fb3064f267047d300b76f3":"Nine Node",
            "0xf26f8e060b241ad3ff9c5821d85544b21419ef2730c5c64dc7bd89bfc3710885":"Tokencan.net Node",
            "0x6cd71031cdc64db40816eecb01f65b3b56cdbde843bcd0797424e6184bf88e07":"Magic Node",
            "0xe058df83ef89df930aca3fed0c949ab6de1ccb51d2a87fc056631f11e3cd2ccc":"Water Separation",
            "0x24a0a75882b3be4eb1b5ceee17c6f253d147cfbc43cb35e798cd5560f43d93ab":"QBTC NODE",
            "0x5c7dbbe680d7687e23b0dd927c37e77f051d701d57ef118f2e7ef5ad9bf50c51":"SERO Cat Node",
            "0x36e40dc58bf2db57f1831a02c2a248b0fa378ba3a99908c17a63ed622174d496":"SERO Monkey Node",
            "0x16a800063ef8345f8441de438d97c8a49bc3373778e45a2f2cd25699be8d2c2b":"MXC Node",
            "0xfd703f582147e4725a9b3567635dde7e442182c1f28e315360e07692a968183b":"Super Node",
            "0xb688861cf4b0f74509a4a4a0bdcdbe80ec0b5647bfcc41973498c9b8152fce41":"OKEx Pool",
            "0xce85c13d0942c3118873facf99669d193711bdca77b3de36668fdf29429bd8e5":"Sun Chase Node",
            "0x92af05b6d6c27a73eb8489f100c859cf3167298915cbae7262345e17aee6aa2c":"Bilaxy Node",
            "0x7fbbc2f11ac02255b48c8ab534d08a0e2fa134d483794ca9294ad296f8e71df7":"HCOIN",
            "0x3800ad263e7b866d5fb21a59f867e415e6d8515b9ae58c84a9c1cdb38e7ba729":"BiKi Pool",
        }
    },

    init: function () {
        var that = this;

        that.getAccountlist();

        $('.close').bind('click', function () {
            $('.modal-deploy').hide();
            $('.modal-footer button:eq(1)').unbind('click');
            $("#sub1").attr('disabled', false);
        });

        $('.register').bind('click', function () {
            window.location.href = 'stake-register.html';
        });

        $('.showShareDetail').bind('click', function () {
            window.location.href = 'stake-detail.html';
        });

        setInterval(function () {
            that.getAccountlist();
        },10000);

        setTimeout(function () {
            that.stakeList();
        },50)

        // setTimeout(function () {
        //     $('.buyShare').bind('click', function () {
        //         var poolId = $(this).attr('attpoolid');
        //         window.location.href = 'stake-buy.html?id=' + poolId;
        //     });
        // },1000)

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
                $('.navbar-nav li:eq(3) a').text($.i18n.prop('navbar_dapps'));

                $('h3:eq(0)').text($.i18n.prop('stake_share_title'));
                $('.showShareDetail').text($.i18n.prop('stake_share_button'));

                $('.totalProfit div:eq(0)').text($.i18n.prop('stake_share_profit'));
                $('.totalAmount div:eq(0)').text($.i18n.prop('stake_share_totalAmount'));
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
                $('thead tr td:eq(3)').text($.i18n.prop('stake_pool_voted_node'));
                $('thead tr td:eq(4)').text($.i18n.prop('stake_pool_voted_solo'));
                $('thead tr td:eq(5)').text($.i18n.prop('stake_pool_missed'));
                $('thead tr td:eq(6)').text($.i18n.prop('stake_pool_missed_rate'));
                $('thead tr td:eq(7)').text($.i18n.prop('stake_pool_fee'));
                $('thead tr td:eq(8)').text($.i18n.prop('stake_pool_shareNum'));
                $('thead tr td:eq(9)').text($.i18n.prop('stake_pool_lastpaytime'));
                $('thead tr td:eq(10)').text($.i18n.prop('stake_pool_operation'));
                $('.buyShare').text($.i18n.prop('stake_pool_buyShare'));

                $('#myModal ul li:eq(0) div div:eq(0)').text($.i18n.prop('send_tx_from'));
                $('#myModal ul li:eq(1) div div:eq(0)').text($.i18n.prop('stake_pool_id'));
                $('#myModal ul li:eq(2) div div:eq(0)').text($.i18n.prop('send_tx_pwd'));
                $('#myModal ul li:eq(3) div div:eq(0)').text($.i18n.prop('stake_pool_close_note'));

                $('#myModal button:eq(1)').text($.i18n.prop('send_tx_cancel'));
                $('#myModal button:eq(2)').text($.i18n.prop('send_tx_confirm'));

                $('#modifyModal h5').text($.i18n.prop('stake_pool_modify'));
                $('#modifyModal label:eq(0)').text($.i18n.prop('stake_register_from'));
                $('#modifyModal label:eq(1)').text($.i18n.prop('stake_register_address'));
                $('#modifyModal label:eq(2)').text($.i18n.prop('stake_register_fee'));
                $('#modifyModal label:eq(3)').text($.i18n.prop('stake_register_password'));

                $('#modifyModal button:eq(1)').text($.i18n.prop('send_tx_cancel'));
                $('#modifyModal button:eq(2)').text($.i18n.prop('send_tx_confirm'));

                $('.closeNode').text($.i18n.prop('stake_pool_close'));
                $('.modifyNode').text($.i18n.prop('stake_pool_modify'));

            }
        });
    },

    closeStakePool: function(pk,id,idPkr){
        $('#myModal ul li:eq(0) div div:eq(1)').text(pk);
        $('#myModal ul li:eq(1) div div:eq(1)').text(id);
        $('#myModal').modal({backdrop: 'static', keyboard: false});

        $('#myModal button:eq(2)').unbind('click').bind('click',function(){
            $(this).attr('disabled',true).text($.i18n.prop('send_tx_sending'));
            var password = $('#password').val();
            if (!password){
                $('#toast2 div').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('stake_pool_password_place'));
                $('#toast2').toast('show');
                $(this).attr('disabled',false).text($.i18n.prop('send_tx_confirm'));
            }else{
                var biz = {};
                biz.From = pk;
                biz.IdPkr = idPkr;
                biz.Password=password;

                var that = this;
                Common.postAsync('stake/close',biz,{},function (res) {
                    if(res.base.code === 'SUCCESS'){
                        $(that).attr('disabled',false).text($.i18n.prop('send_tx_confirm'));
                        $('#toast2 div').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                        $('#toast2').toast('show');
                        setTimeout(function () {
                            window.location.href = 'account-detail.html?pk='+pk;
                        }, 1500);
                    }else{
                        $('#toast2 div').removeClass('alert-success').addClass('alert-danger').text(Common.convertErrors(res.base.desc));
                        $('#toast2').toast('show');
                        $(that).attr('disabled',false).text($.i18n.prop('send_tx_confirm'));
                    }

                })
            }
        });
    },

    modifyStakePool: function(pk,id,voteAddress,feeRate,idPkr){
       $('#from').empty().append(`
            <option value=${pk}>${pk}</option>
       `);

        $('#idPkr').val(idPkr);
        $('#voteAddress').val(voteAddress);
        $('#feeRate').val(feeRate);
        $('#modifyModal').modal({backdrop: 'static', keyboard: false});
    },

    modifyStakeConfirm:function(){
        var from = $("#from").val();
        var vote = $("#voteAddress").val();
        var idPkr = $('#idPkr').val();
        var feeRate = $("#feeRate").val();

        var password = $('#passwordModify').val();
        if (!password){
            $('#toastModifyStake div').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('stake_pool_password_place'));
            $('#toastModifyStake').toast('show');
            $(this).attr('disabled',false).text($.i18n.prop('send_tx_confirm'));
        }else {
            var biz = {
                From: from,
                Vote: vote,
                FeeRate: new BigNumber(feeRate).multipliedBy(100).toString(10),
                Type: "modify",
                IdPkr: idPkr,
                Password:password,
            }
            Common.postAsync('stake/register', biz, {}, function (res) {
                if (res.base.code === 'SUCCESS') {
                    $('#toastModifyStake div').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                    $('#toastModifyStake').toast('show');

                    $('#toastModifyStake button:eq(2)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                    $('#sub1').attr('disabled', false);
                    setTimeout(function () {
                        window.location.href = 'account-detail.html?pk='+from;
                    }, 1500);
                } else {
                    $('#toastModifyStake div').removeClass('alert-success').addClass('alert-danger').text(Common.convertErrors(res.base.desc));
                    $('#toastModifyStake').toast('show');
                    $('#toastModifyStake button:eq(2)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                }
                $('#sub1').attr('disabled', false);
            })
        }
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
                    var totalAmount = new BigNumber(0);

                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        var acName = "Account"+(i + 1);
                        if (data.Name){
                            acName = data.Name;
                        }
                        that.account[data.MainPKr]= acName +"("+data.PK.substring(0, 8) + " ... " + data.PK.substring(data.PK.length - 8, data.PK.length)+")";
                        that.account[data.MainOldPKr]= acName +"("+data.PK.substring(0, 8) + " ... " + data.PK.substring(data.PK.length - 8, data.PK.length)+")";
                        that.accountMap[data.MainPKr]= data;
                        that.accountMap[data.MainOldPKr]= data;


                        Common.post('share/my', data.MainPKr, {}, function (res2) {
                            if (res2.biz && res2.base.code === 'SUCCESS') {
                                if (res2.biz && res2.biz.length > 0) {
                                    var data = res2.biz[0];

                                    totalAmount = totalAmount.plus(new BigNumber(data.totalAmount?data.totalAmount:"0x0",16));

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

                                    $('.totalProfit span:eq(1)').text(totalProfit.toFixed(6) + ' SERO');
                                    $('.totalAmount span:eq(1)').text(totalAmount.dividedBy(Common.baseDecimal).toFixed(6) + ' SERO');
                                    $('.expireShareNum span:eq(1)').text(expireShareNum.toString(10));
                                    $('.leftShareNum span:eq(1)').text(leftShareNum.toString(10));
                                    $('.missShareNum span:eq(1)').text(missShareNum.toString(10));
                                    $('.hasShareNum span:eq(1)').text(totalShareNum.minus(expireShareNum).minus(missShareNum).minus(leftShareNum).toString(10));
                                    $('.totalShareNum span:eq(1)').text(totalShareNum.toString(10));
                                }
                            };
                        });

                        Common.post('share/my', data.MainOldPKr, {}, function (res2) {
                            if (res2.biz && res2.base.code === 'SUCCESS') {
                                if (res2.biz && res2.biz.length > 0) {
                                    var data = res2.biz[0];
                                    totalAmount = totalAmount.plus(new BigNumber(data.totalAmount?data.totalAmount:"0x0",16));
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

                                    $('.totalProfit span:eq(1)').text(totalProfit.toFixed(6) + ' SERO');
                                    $('.totalAmount span:eq(1)').text(totalAmount.dividedBy(Common.baseDecimal).toFixed(6) + ' SERO');
                                    $('.expireShareNum span:eq(1)').text(expireShareNum.toString(10));
                                    $('.leftShareNum span:eq(1)').text(leftShareNum.toString(10));
                                    $('.missShareNum span:eq(1)').text(missShareNum.toString(10));
                                    $('.hasShareNum span:eq(1)').text(totalShareNum.minus(expireShareNum).minus(missShareNum).minus(leftShareNum).toString(10));
                                    $('.totalShareNum span:eq(1)').text(totalShareNum.toString(10));
                                }
                            }
                        });

                    }
                }
            }
        })
    },

    stakeList: function () {
        var that = this;
        $('tbody').empty();

        Common.post('stake', {}, {}, function (res) {
            if (res.base.code === 'SUCCESS') {
                var lan = $.cookie('language')?$.cookie('language'):"en_US";
                var _stakeName = that.stakeName[lan];
                var dataArray = res.biz;
                for (var data of dataArray) {

                    var isMy = `<span class="text-primary">${that.account[data.idPkr]?"Created by: "+that.account[data.idPkr]:""}</span><br/>`;

                    var state = `<span class="text-success">OPENING</span>`;

                    var buyButton = `<button class="btn btn-outline-primary btn-block btn-sm buyShare" attpoolid="${data.id}" onclick="goBuy(${"'"+data.id+"'"})">${$.i18n.prop('stake_pool_buyShare')}</button>`;
                    var closeButton = "";
                    var modifyButton = "";

                    if(that.account[data.idPkr] && !data.closed){
                        closeButton = `<button class="btn btn-outline-danger btn-block btn-sm closeNode" onclick="closeStake(${"'"+that.accountMap[data.idPkr].PK+"','"+data.id+"','"+data.idPkr+"'"})">${$.i18n.prop('stake_pool_close')}</button>`;
                        modifyButton = `<button class="btn btn-outline-info btn-block btn-sm modifyNode" onclick="modifyStake(${"'"+that.accountMap[data.idPkr].PK+"','"+data.id+"','"+data.voteAddress+"','"+(new BigNumber(data.fee?data.fee:"0x0", 16).div(100).toFixed(2))+"','"+data.idPkr+"'"})">${$.i18n.prop('stake_pool_modify')}</button>`;
                    }
                    if (data.closed){
                        buyButton = "";
                        modifyButton = "";
                        closeButton = "";

                        state = `<span class="text-danger">CLOSED</span>`;
                    }
                    var choiceNum = new BigNumber(data.choicedNum?data.choicedNum:"0x0", 16);
                    var missed = new BigNumber(data.missedNum?data.missedNum:"0x0", 16);
                    var wishVoteNum = new BigNumber(data.wishVoteNum?data.wishVoteNum:"0x0", 16);
                    var nodeVoted = choiceNum.minus(missed);

                    var soloVoted = missed.minus(wishVoteNum);
                    var missRate = "--";
                    if (nodeVoted.comparedTo(0)>0){
                        missRate = wishVoteNum.dividedBy(nodeVoted).multipliedBy(100).toFixed(2)+"%";
                    }

                    var profit =  `<span class="text-success">${new BigNumber(data.profit?data.profit:"0x0", 16).dividedBy(Common.baseDecimal).toFixed(6)}</span>`;

                    $('tbody').append(`
                    <tr>
                        <td class="text-break">${data.id}</td>
                        <td class="text-break">
                            <span class="text-info">${_stakeName[data.id]?_stakeName[data.id]:""}</span><br/>
                            ${data.own.substring(0,8) + " ... " + data.own.substring(data.own.length-8)}<br/>
                            ${isMy}
                            ${that.account[data.own]?"Profit: "+profit:""}
                        </td>
                        <td>${state}</td>
                        <td>${nodeVoted.toString(10)}</td>
                        <td>${soloVoted.toString(10)}</td>
                        <td>${wishVoteNum.toString(10)}</td>
                        <td><span class="text-danger">${missRate}</span> </td>
                        <td>${new BigNumber(data.fee?data.fee:"0x0", 16).div(100).toFixed(2)}%</td>
                        <td>${new BigNumber(data.shareNum?data.shareNum:"0x0", 16).toString()}</td>
                        <td>${new BigNumber(data.lastPayTime?data.lastPayTime:"0x0", 16).toString(10)}</td>
                        <td>${buyButton}${modifyButton}${closeButton}</td>
                    </tr>
               `);
                }
                $('#dataTable').DataTable();
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
                $('.navbar-nav li:eq(3) a').text($.i18n.prop('navbar_dapps'));

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
                        var acName = "Account"+(i + 1);
                        if (data.Name){
                            acName = data.Name;
                        }
                        if (data.Balance) {
                            var balanceObj = data.Balance;
                            for (var currency of Object.keys(balanceObj)) {
                                if (currency === 'SERO') {
                                    balance = new BigNumber(balanceObj[currency]).dividedBy(Common.baseDecimal).toFixed(6);
                                    $('.address').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${acName +": "+data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8) }  ${ balance + ' ' + currency}</option>`);
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
        $('#myModal').modal({backdrop: 'static', keyboard: false});

        $('.modal-footer button:eq(1)').bind('click', function () {
            var password = $('#password').val();
            if (!password){
                $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('stake_pool_password_place'));
                $('.toast').toast('show');
                $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
            }else{
                if(vote === ''){
                    $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('stake_register_address_tips'));
                    $('.toast').toast('show');
                }else{
                    $('.modal-footer button:eq(1)').attr('disabled', true).text($.i18n.prop('send_tx_sending'));
                    var biz = {
                        From: from,
                        Vote: vote,
                        FeeRate: new BigNumber(feeRate).multipliedBy(100).toString(10),
                        Password:password,
                    }
                    Common.postAsync('stake/register', biz, {}, function (res) {
                        if (res.base.code === 'SUCCESS') {
                            $('.toast-body').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                            $('.toast').toast('show');

                            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                            $('#sub1').attr('disabled', false);
                            setTimeout(function () {
                                window.location.href = 'account-detail.html?pk='+from;
                            }, 1500);
                        } else {
                            $('.toast-body').removeClass('alert-success').addClass('alert-danger').text(Common.convertErrors(res.base.desc));
                            $('.toast').toast('show');
                            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                        }
                    })
                }
            }


        });
    }


};

var StakeBuyer = {

    poolId: '',
    mainPKr:{},

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
        $('.toast').toast({animation: true, autohide: true, delay: 2000});

        $('#amount').bind('input',function () {
            that.estimateShares();
        })
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

                $('.breadcrumb li:eq(0) a').text($.i18n.prop('stake_pool_title'));
                $('.breadcrumb li:eq(1)').text($.i18n.prop('stake_pool_buyShare'));
                $('.form-group:eq(0) label').text($.i18n.prop('share_buy_poolId'));
                $('.form-group:eq(1) label').text($.i18n.prop('share_buy_from'));
                // $('.form-group:eq(2) label').text($.i18n.prop('share_buy_address'));
                // $('.form-group:eq(2) .invalid-feedback').text($.i18n.prop('share_buy_address_tips'));
                $('.form-group:eq(2) label').text($.i18n.prop('share_buy_amount'));
                $('.form-group:eq(2) .invalid-feedback').text($.i18n.prop('share_buy_amount_tips'));
                $('#address').attr('placeholder', $.i18n.prop('share_buy_address_tips'));
                $('#amount').attr('placeholder', $.i18n.prop('share_buy_amount_tips'));
                $('h4').text($.i18n.prop('stake_pool_buyShare'));
                $('.modal-title').text($.i18n.prop('stake_register_confirm_title'));
                $('.modal-body ul li:eq(0) div div:eq(0)').text($.i18n.prop('share_buy_poolId'));
                $('.modal-body ul li:eq(1) div div:eq(0)').text($.i18n.prop('share_buy_from'));
                // $('.modal-body ul li:eq(2) div div:eq(0)').text($.i18n.prop('share_buy_address'));
                $('.modal-body ul li:eq(2) div div:eq(0)').text($.i18n.prop('share_buy_amount'));
                $('.modal-body ul li:eq(3) div div:eq(0)').text($.i18n.prop('stake_register_password'));
                $('#password').attr('placeholder', $.i18n.prop('stake_register_password_place'));
                $('#sub1').text($.i18n.prop('stake_register_next'));

                $('.modal-footer button:eq(0)').text($.i18n.prop('stake_register_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('stake_register_confirm'));

                $('.estimateShares span:eq(0)').text($.i18n.prop('share_buy_estimate_price'));
                $('.estimateShares span:eq(1)').text($.i18n.prop('share_buy_estimate_total'));

                $('.amount_warning').text($.i18n.prop('share_buy_amount_waring'));


            }
        });
    },

    estimateShares:function () {

        var that = this;
        var from = $(".address").val();
        // var vote = $("#address").val();
        var amount = $("#amount").val();

        var params = {
            from:from,
            vote:that.mainPKr[from],
            value:"0x"+new BigNumber(amount).multipliedBy(Common.baseDecimal).toString(16),
        }

        if (from !== "" ){
            Common.postSeroRpc("stake_estimateShares",[params],function (res) {
                if(res.result){
                    var result = res.result;
                    var avPrice = result.avPrice;
                    // var basePrice = result.basePrice;
                    // "Average Price: "+new BigNumber(avPrice,16).dividedBy(Common.baseDecimal).toFixed(6) + " SERO"
                    var total = result.total;

                    $('.estimateShares strong:eq(0)').text(new BigNumber(avPrice,16).dividedBy(Common.baseDecimal).toFixed(6));
                    $('.estimateShares strong:eq(1)').text(new BigNumber(total,16).toString(10));
                }
            });
        }

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
                        that.mainPKr[data.PK] = data.MainPKr;
                        var balance = new BigNumber(0).toFixed(6);
                        var acName = "Account"+(i + 1);
                        if (data.Name){
                            acName = data.Name;
                        }
                        if (data.Balance) {
                            var balanceObj = data.Balance;
                            for (var currency of Object.keys(balanceObj)) {
                                if (currency === 'SERO') {
                                    balance = new BigNumber(balanceObj[currency]).dividedBy(Common.baseDecimal).toFixed(6);
                                    $('.address').append(`<option value="${data.PK}" ${i === 0 ? 'selected' : ''}>${acName + ": "+data.PK.substring(0, 8) + ' ... ' + data.PK.substring(data.PK.length - 8) }  ${ balance + ' ' + currency}</option>`);
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
        // $('ul:eq(1) li:eq(2) div div:eq(1)').text(vote);
        $('ul:eq(1) li:eq(2) div div:eq(1)').text(amount);
        $('#myModal').modal({backdrop: 'static', keyboard: false});

        $('.modal-footer button:eq(1)').bind('click', function () {
            var password = $('#password').val();
            if (!password){
                $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('stake_pool_password_place'));
                $('.toast').toast('show');
                $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
            }else{
                var estimateShares = parseInt($('.estimateShares strong:eq(1)').text());

                if(estimateShares === 0){
                    $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('share_buy_amount_fail'));
                    $('.toast').toast('show');
                }else{
                    $('.modal-footer button:eq(1)').attr('disabled', true).text($.i18n.prop('send_tx_sending'));
                    var biz = {
                        From: from,
                        Vote: vote,
                        Amount: new BigNumber(amount).multipliedBy(Common.baseDecimal).toString(10),
                        Pool: poolId,
                        GasPrice: new BigNumber(1000000000).toString(10),
                        Password:password,
                    }
                    Common.postAsync('stake/buyShare', biz, {}, function (res) {
                        if (res.base.code === 'SUCCESS') {
                            $('.toast-body').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('send_tx_success'));
                            $('.toast').toast('show');
                            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                            $('#sub1').attr('disabled', false);
                            setTimeout(function () {
                                window.location.href = 'account-detail.html?pk='+from;
                            }, 1500);
                        } else {
                            $('.toast-body').removeClass('alert-success').addClass('alert-danger').text(Common.convertErrors(res.base.desc));
                            $('.toast').toast('show');
                            $('.modal-footer button:eq(1)').attr('disabled', false).text($.i18n.prop('send_tx_confirm'));
                        }
                    })
                }
            }

        });
    }


};

var StakeDetail = {


    init: function () {
        var that = this;

        setTimeout(function () {
            that.getAccountlist();
        },10)


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
                $('.shareDetail').text($.i18n.prop('share_detail_title'));
                $('thead tr td:eq(0) strong').text($.i18n.prop('share_detail_shareId'));
                $('thead tr td:eq(1) strong').text($.i18n.prop('share_detail_poolId'));
                $('thead tr td:eq(2) strong').text($.i18n.prop('share_detail_address'));
                $('thead tr td:eq(3) strong').text($.i18n.prop('share_detail_price'));
                $('thead tr td:eq(4) strong').text($.i18n.prop('share_detail_fee'));
                $('thead tr td:eq(5) strong').text($.i18n.prop('share_detail_profit'));
                $('thead tr td:eq(6) strong').text($.i18n.prop('share_detail_returned'));
                $('thead tr td:eq(7) strong').text($.i18n.prop('share_detail_shared'));
                $('thead tr td:eq(8) strong').text($.i18n.prop('share_detail_time'));

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

                    var avgPrice  = new BigNumber(0);
                    var totalProfit = new BigNumber(0);
                    var totalRemaining = new BigNumber(0);
                    var totalVoted = new BigNumber(0);
                    var totalExpired = new BigNumber(0);
                    var totalMissed = new BigNumber(0);
                    var totalShares = new BigNumber(0);
                    var totalReturnedProfit = new BigNumber(0);
                    var totalReturnedNumber = new BigNumber(0);
                    var totalAmount = new BigNumber(0);

                    var count = 0 ;
                    for (var i = 0; i < dataArray.length; i++) {
                        var data = dataArray[i];
                        Common.post('share/my', data.MainPKr, {}, function (res) {
                            if (res.base.code === 'SUCCESS') {
                                if (res.biz.length > 0) {
                                    var dataShare = res.biz[0];
                                    var shareIds = dataShare.shareIds;
                                    totalAmount = totalAmount.plus(new BigNumber(dataShare.totalAmount?dataShare.totalAmount:"0x0",16));
                                    for (let shareId of shareIds) {
                                        Common.post('stake/getShare', shareId, {}, function (res) {
                                            var share = res.biz;
                                            if (res.base.code === 'SUCCESS') {
                                                count++;
                                                var voted = new BigNumber(share.total, 16).minus(new BigNumber(share.remaining?share.remaining:"0x0", 16)).minus(new BigNumber(share.missed?share.missed:"0x0", 16)).minus(new BigNumber(share.expired?share.expired:"0x0", 16));

                                                avgPrice = avgPrice.plus(new BigNumber(share.price, 16));

                                                totalProfit = totalProfit.plus(new BigNumber(share.profit, 16));
                                                totalRemaining = totalRemaining.plus(new BigNumber(share.remaining?share.remaining:"0x0", 16));
                                                totalVoted = totalVoted.plus(voted);
                                                totalExpired = totalExpired.plus(new BigNumber(share.expired?share.expired:"0x0", 16));
                                                totalMissed = totalMissed.plus(new BigNumber(share.missed?share.missed:"0x0", 16));
                                                totalShares = totalShares.plus(new BigNumber(share.total,16));
                                                totalReturnedProfit = totalReturnedProfit.plus(new BigNumber(share.returnProfit?share.returnProfit:"0x0",16));
                                                totalReturnedNumber = totalReturnedNumber.plus(new BigNumber(share.returnNum?share.returnNum:"0x0",16));

                                                var acName = "Account"+(i+1);
                                                if (data.Name){
                                                    acName = data.Name;
                                                }
                                                $('tbody').append(`
                                                <tr>
                                                <td class="text-break">${share.id}</td>
                                                <td class="text-break">${share.pool}</td>
                                                <td class="text-break">
                                                    ${acName}(${data.PK.substring(0, 5) + " ... " + data.PK.substring(data.PK.length - 5)})<br/><br/>
                                                    Share Address: ${dataShare.addr}
                                                </td>
                                                <td>${new BigNumber(share.price, 16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                                <td>${(parseFloat(new BigNumber(share.fee,16).toString(10)) / 100).toFixed(2)}%</td>
                                                <td>${new BigNumber(share.profit, 16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                                <td>
                                                    <strong class="text-info">${new BigNumber(share.returnProfit?share.returnProfit:"0x0", 16).dividedBy(Common.baseDecimal).toFixed(6)}</strong>
                                                </td>
                                                <td>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_remaining')}: </small><strong class="text-info">${new BigNumber(share.remaining?share.remaining:"0x0", 16).toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_voted')}: </small><strong class="text-info">${voted.toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_expired')}: </small><strong class="text-info">${new BigNumber(share.expired?share.expired:"0x0", 16).toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_missed')}: </small><strong class="text-info">${new BigNumber(share.missed?share.missed:"0x0", 16).toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_total')}: </small><strong class="text-info">${new BigNumber(share.total,16).toString(10)}</strong><br/><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_number')}: </small><strong class="text-danger">${new BigNumber(share.returnNum?share.returnNum:"0x0", 16).toString(10)}</strong><br/>
                                                </td>
                                                <td>${toTime(share.timestamp)}</td>
                                                </tr>
                                            `);
                                            }
                                        });
                                    }
                                }
                            }
                        });
                        Common.post('share/my', data.MainOldPKr, {}, function (res) {
                            if (res.base.code === 'SUCCESS') {
                                if (res.biz.length > 0) {
                                    var dataShare = res.biz[0];
                                    var shareIds = dataShare.shareIds;
                                    totalAmount = totalAmount.plus(new BigNumber(dataShare.totalAmount?dataShare.totalAmount:"0x0",16));
                                    for (let shareId of shareIds) {
                                        Common.post('stake/getShare', shareId, {}, function (res) {
                                            var share = res.biz;
                                            if (res.base.code === 'SUCCESS') {
                                                count++;
                                                var voted = new BigNumber(share.total, 16).minus(new BigNumber(share.remaining?share.remaining:"0x0", 16)).minus(new BigNumber(share.missed?share.missed:"0x0", 16)).minus(new BigNumber(share.expired?share.expired:"0x0", 16));

                                                avgPrice = avgPrice.plus(new BigNumber(share.price, 16));
                                                totalProfit = totalProfit.plus(new BigNumber(share.profit, 16));
                                                totalRemaining = totalRemaining.plus(new BigNumber(share.remaining?share.remaining:"0x0", 16));
                                                totalVoted = totalVoted.plus(voted);
                                                totalExpired = totalExpired.plus(new BigNumber(share.expired?share.expired:"0x0", 16));
                                                totalMissed = totalMissed.plus(new BigNumber(share.missed?share.missed:"0x0", 16));
                                                totalShares = totalShares.plus(new BigNumber(share.total,16));
                                                totalReturnedProfit = totalReturnedProfit.plus(new BigNumber(share.returnProfit,16));
                                                totalReturnedNumber = totalReturnedNumber.plus(new BigNumber(share.returnNum,16));

                                                var acName = "Account"+(i+1);
                                                if (data.Name){
                                                    acName = data.Name;
                                                }
                                                $('tbody').append(`
                                                <tr>
                                                <td class="text-break">${share.id}</td>
                                                <td class="text-break">${share.pool}</td>
                                                <td class="text-break">
                                                    ${acName}(${data.PK.substring(0, 5) + " ... " + data.PK.substring(data.PK.length - 5)})<br/><br/>
                                                    Share Address: ${dataShare.addr}
                                                </td>
                                                <td>${new BigNumber(share.price, 16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                                <td>${(parseFloat(new BigNumber(share.fee,16).toString(10)) / 100).toFixed(2)}%</td>
                                                <td>${new BigNumber(share.profit, 16).dividedBy(Common.baseDecimal).toFixed(6)}</td>
                                                <td>
                                                    <strong class="text-info">${new BigNumber(share.returnProfit?share.returnProfit:"0x0", 16).dividedBy(Common.baseDecimal).toFixed(6)}</strong>
                                                </td>
                                                <td>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_remaining')}: </small><strong class="text-info">${new BigNumber(share.remaining?share.remaining:"0x0", 16).toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_voted')}: </small><strong class="text-info">${voted.toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_expired')}: </small><strong class="text-info">${new BigNumber(share.expired?share.expired:"0x0", 16).toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_missed')}: </small><strong class="text-info">${new BigNumber(share.missed?share.missed:"0x0", 16).toString(10)}</strong><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_total')}: </small><strong class="text-info">${new BigNumber(share.total,16).toString(10)}</strong><br/><br/>
                                                    <small class="text-gray-500">${$.i18n.prop('share_detail_number')}: </small><strong class="text-danger">${new BigNumber(share.returnNum?share.returnNum:"0x0", 16).toString(10)}</strong><br/>
                                                </td>
                                                <td>${toTime(share.timestamp)}</td>
                                                </tr>
                                            `);
                                            }
                                        });
                                    }
                                }
                            }
                        });
                    }

                    $('tfoot tr td:eq(3) strong').text("Average "+avgPrice.dividedBy(Common.baseDecimal).dividedBy(count===0?1:count).toFixed(6));
                    $('tfoot tr td:eq(5) strong').text(totalProfit.dividedBy(Common.baseDecimal).toFixed(6));
                    $('tfoot tr td:eq(6) strong').append(
                        `
                        <span class="text-info">${totalReturnedProfit.dividedBy(Common.baseDecimal).toFixed(6)}</span>                        
                        `
                    );

                    $('tfoot tr td:eq(7) strong').append(`
                        <small class="text-gray-500">${$.i18n.prop('share_detail_remaining')}: </small><span class="text-info">${totalRemaining.toString(10)}</span><br/>
                        <small class="text-gray-500">${$.i18n.prop('share_detail_voted')}: </small><span class="text-info">${totalVoted.toString(10)}</span><br/>
                        <small class="text-gray-500">${$.i18n.prop('share_detail_expired')}: </small><span class="text-info">${totalExpired.toString(10)}</span><br/>
                        <small class="text-gray-500">${$.i18n.prop('share_detail_missed')}: </small><span class="text-info">${totalMissed.toString(10)}</span><br/>
                        <small class="text-gray-500">${$.i18n.prop('share_detail_total')}: </small><span class="text-info">${totalShares.toString(10)}</span><br/><br/>
                        <small class="text-gray-500">${$.i18n.prop('share_detail_number')}: </small><span class="text-danger">${totalReturnedNumber.toString(10)}</span><br/>
                        <small class="text-gray-500">${$.i18n.prop('share_detail_totalAmount')}: </small><span class="text-success">${totalAmount.dividedBy(Common.baseDecimal).toFixed(6)}</span><br/>
                    `);
                }
                $('#dataTable').DataTable();
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

function goBuy (poolId) {
    window.location.href = 'stake-buy.html?id=' + poolId;
}

function toTime(timstamp) {

    if(timstamp){
        var tt = parseFloat(new BigNumber(timstamp,16).toString(10));
        var dt = new Date(tt*1000);
        // return dt.toJSON().substr(0, 19).replace('T', ' ').replace(/-/g, '.');
        return convertUTCTimeToLocalTime(dt);
    }else{
        return "";
    }
}

function convertUTCTimeToLocalTime (UTCDateString) {
    if(!UTCDateString){
        return '-';
    }
    function formatFunc(str) {    //格式化显示
        return str > 9 ? str : '0' + str
    }
    var date2 = new Date(UTCDateString);     //这步是关键
    var year = date2.getFullYear();
    var mon = formatFunc(date2.getMonth() + 1);
    var day = formatFunc(date2.getDate());
    var hour = date2.getHours();
    var noon = hour >= 12 ? 'PM' : 'AM';
    hour = hour>=12?hour-12:hour;
    hour = formatFunc(hour);
    var min = formatFunc(date2.getMinutes());
    var dateStr = year+'-'+mon+'-'+day+' '+noon +' '+hour+':'+min;
    return dateStr;
}