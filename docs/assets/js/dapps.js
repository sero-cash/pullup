var DApps = {

    init: function () {
        var that = this;
        setTimeout(function () {
            that.genPageData();
        }, 10)

    },

    goBrowser(url){
        window.location.href = "./browser.html?dapp_url="+encodeURIComponent(url);
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
                $('.navbar-nav li:eq(0) a span').text($.i18n.prop('navbar_home'));
                $('.navbar-nav li:eq(1) a span').text($.i18n.prop('navbar_send'));
                $('.navbar-nav li:eq(2) a span').text($.i18n.prop('navbar_stake'));
                $('.navbar-nav li:eq(3) a span').text($.i18n.prop('navbar_dapps'));

                $('#myModal .modal-title').text($.i18n.prop('dapps_modal_title'));
                $('#myModal .modal-body p').text($.i18n.prop('dapps_modal_body'));
                $('#myModal .modal-footer button:eq(0)').text($.i18n.prop('dapps_button_cancel'));
                $('#myModal .modal-footer button:eq(1)').text($.i18n.prop('dapps_button_enter'));

                $('.add-dapp h3').text($.i18n.prop('navbar_dapps'));
                $('.add-dapp i').text($.i18n.prop('dapps_modal_add'));

                $('#addDappModal .modal-footer button:eq(0)').text($.i18n.prop('dapp_token_modal_button_cancel'));
                $('#addDappModal .modal-footer button:eq(1)').text($.i18n.prop('dapps_modal_add'));
                $('#addDappModal .modal-title').text($.i18n.prop('dapps_modal_add'));
            }
        });
        $('.toast').toast({animation: true, autohide: true, delay: 2000})
        this.genPageData();
    },

    dapps_en_US: [
        {
            img: "./assets/img/token.png",
            title: "SRC20 Token Tracker",
            desc: "SRC20 is a standard interface for anonymous tokens. This standard provides basic functionality to transfer tokens.",
            author: "sero.cash",
            url: "./views/contract/token.html",
            showTips: false,
            state: 1,
        },
        {
            img: "./assets/img/asnow.jpeg",
            title: "ASNOW",
            desc: "",
            author: "asnow.com",
            url: "http://asnowhero.gitee.io/asnow-pullup/",
            showTips: true,
            state: 1,
        },
        {
            "img": "https://sero-cash.github.io/voter/vote.jpg",
            "title": "Vote",
            "desc": "Support or Oppose, Who can win?",
            "author": "nobody",
            "url": "https://sero-cash.github.io/voter/v1/",
            showTips: true,
            "state": 1
        },
        {
            "img": "http://osaifu.io:3000/logo.png",
            "title": "LUCKY",
            "desc": "Decentralized Lucky Game",
            "author": "Satoshi Nakamoto",
            "url": "http://osaifu.io:3000/",
            showTips: true,
            "state": 1
        },
        {
            img: "./assets/img/sanguo.png",
            title: "超零三国-无限穿越",
            desc: "抽卡功能已上线！人物全服唯一，是属于你的数字资产，可交易，可赚钱！ 后续马上推出游历闯关、红颜约会、武将养成、养马、三国杀等丰富玩法！  超零平台第一款真正意义上的游戏",
            author: "盘古",
            url: "http://122.112.190.100/www/slg/slg.html",
            showTips: true,
            state: 1,
        },
        {
            "img": "http://39.98.58.15/logo.png",
            "title": "ACES",
            "desc": "AIPP로 ACES게임에 참여하는것을 환영합니다.",
            "author": "Tokencan",
            "url": "http://39.98.58.15",
            showTips: true,
            "state": 1
        },
        {
            "img": "./assets/img/TableName.png",
            "title": "TableGame",
            "desc": "TableGame是由舍得社区基于SERO开发的智能合约游戏。用户可自行质押SERO创建游戏桌并邀请朋友加入，每局游戏根据大家的幸运数开奖，中奖者可实时提现！祝大家好运！",
            "author": "TableGame",
            "url": "http://58.64.184.143:3000/#/",
            showTips: true,
            "state": 1
        },
        {
            img: "https://fpsc2019.github.io/fpsc-pullup/logo.png",
            title: "HAPY",
            desc: "HAPY는 해피쉐어링 플랫폼내에서 사용할 수 있는 토큰으로서 물품을 구매하거나 투자 상품에 참여하기 위해서 필요한 암호화폐입니다.\n" +
                "해피쉐어링은 글로벌 공유 커뮤니티입니다. 현재 일본을 주축으로 이미 10만명 규모의 유저들이 활동 중이며 더 많은 유저를 확보하기 위해 중국 베트남 등 글로벌 시장 진출을 준비하고 있습니다.",
            author: "HAPY",
            url: "https://fpsc2019.github.io/fpsc-pullup/",
            showTips: true,
            state: 1,
        }
    ],
    dapps_zh_CN: [
        {
            img: "./assets/img/token.png",
            title: "SRC20 Token Tracker",
            desc: "SRC20是匿名Token的标准接口, 该标准提供了转移Token的基本功能。",
            author: "sero.cash",
            url: "./views/contract/token.html",
            showTips: false,
            state: 1,
        },
        {
            img: "./assets/img/asnow.jpeg",
            title: "ASNOW",
            desc: "",
            author: "asnow.com",
            url: "http://asnowhero.gitee.io/asnow-pullup/",
            showTips: true,
            state: 1,
        },{
            "img": "http://sero-cash.gitee.io/voter/vote.jpg",
            "title": "Vote",
            "desc": "Support or Oppose, Who can win?",
            "author": "nobody",
            "url": "http://sero-cash.gitee.io/voter/v1/",
            showTips: true,
            "state": 1
        },
        {
            "img": "http://osaifu.io:3000/logo.png",
            "title": "LUCKY",
            "desc": "Decentralized Lucky Game",
            "author": "Satoshi Nakamoto",
            "url": "http://osaifu.io:3000/",
            showTips: true,
            "state": 1
        },
        {
            img: "./assets/img/sanguo.png",
            title: "超零三国-无限穿越",
            desc: "抽卡功能已上线！人物全服唯一，是属于你的数字资产，可交易，可赚钱！ 后续马上推出游历闯关、红颜约会、武将养成、养马、三国杀等丰富玩法！  超零平台第一款真正意义上的游戏",
            author: "盘古",
            url: "http://122.112.190.100/www/slg/slg.html",
            showTips: true,
            state: 1,
        },
        {
            "img": "http://39.98.58.15/logo.png",
            "title": "ACES",
            "desc": "AIPP로 ACES게임에 참여하는것을 환영합니다.",
            "author": "Tokencan",
            "url": "http://39.98.58.15",
            showTips: true,
            "state": 1
        },
        {
            "img": "./assets/img/TableName.png",
            "title": "TableGame",
            "desc": "TableGame是由舍得社区基于SERO开发的智能合约游戏。用户可自行质押SERO创建游戏桌并邀请朋友加入，每局游戏根据大家的幸运数开奖，中奖者可实时提现！祝大家好运！",
            "author": "TableGame",
            "url": "http://58.64.184.143:3000/#/",
            showTips: true,
            "state": 1
        },
        {
            img: "https://fpsc2019.gitee.io/fpsc-pullup/logo.png",
            title: "HAPY",
            desc: "HAPY는 해피쉐어링 플랫폼내에서 사용할 수 있는 토큰으로서 물품을 구매하거나 투자 상품에 참여하기 위해서 필요한 암호화폐입니다.\n" +
                "해피쉐어링은 글로벌 공유 커뮤니티입니다. 현재 일본을 주축으로 이미 10만명 규모의 유저들이 활동 중이며 더 많은 유저를 확보하기 위해 중국 베트남 등 글로벌 시장 진출을 준비하고 있습니다.",
            author: "HAPY",
            url: "https://fpsc2019.gitee.io/fpsc-pullup/",
            showTips: true,
            state: 1,
        }
    ],

    showAddModal() {
        $('#addDappModal').modal('show');
    },

    addDapp() {
        var that = this;
        var url = $('#url').val();

        url = url.trim();
        var biz = {
            operation: "add",
            url: url
        };

        Common.postAsync('dapp/set', biz, {}, function (res) {
            if (res.base.code === "SUCCESS") {
                that.genPageData();
                $('.toast-body').removeClass('alert-danger').addClass('alert-success').text($.i18n.prop('dapps_modal_add_success'));
                $('.toast').toast('show');
                setTimeout(function () {
                    $('#addDappModal').modal('hide');
                    $('#url').val('');
                },1000);
            } else {
                $('.toast-body').removeClass('alert-success').addClass('alert-danger').text($.i18n.prop('dapps_modal_add_fail'));
                $('.toast').toast('show');
            }
            $('#sub1').attr('disabled', false);
        });

    },

    removeDapp(dappId) {
        var that = this;
        var biz = {
            operation: "remove",
            id: dappId,
        }
        Common.postAsync('dapp/set', biz, {}, function (res) {
            if (res.base.code === "SUCCESS") {
                that.genPageData();
            } else {

            }
        })
    },

    genPageData() {
        var that = this;
        var lang = $.cookie('language');

        if (!lang) {
            lang = "en_US";
            $.cookie('language', lang);
        }
        var data = [];
        if (lang === "zh_CN") {
            data = that.dapps_zh_CN;
        } else if (lang === "en_US") {
            data = that.dapps_en_US
        }

        $(".dapp-data").empty();
        for (var i = 0; i < data.length; i++) {
            var dapp = data[i];
            if (dapp.state === 1) {
                $('.dapp-data').append(`
                    <div class="col-lg-3 col-md-4 col-sm-6" style="padding-bottom:15px;">
                        <div class="card shadow">
                            <img src="${dapp.img}" with="390" style="max-height:280px" class="card-img-top">
                            <div class="card-body" style="height:130px;">
                                <h6 class="card-title text-dark">${dapp.title}</h6>
                                <p class="card-text">${dapp.desc}</p>
                            </div>
                            <div class="card-footer text-right">
                                <a href="${dapp.showTips ? "#" : dapp.url}" class="btn btn-sm btn-primary dapp-btn text-uppercase" dapp-name="${dapp.title}" dapp-url="${dapp.url}">${$.i18n.prop('dapps_button_enter')}</a>
                            </div>
                        </div>
                    </div>
                `);
            } else if (dapp.state === 0) {
                $('.dapp-data').append(`
                    <div class="col-lg-3 col-md-4 col-sm-6" style="padding-bottom:15px;">
                        <div class="card shadow">
                            <img src="${dapp.img}" with="390" style="max-height:280px"  class="card-img-top">
                            <div class="card-body" style="height:130px;">
                                <h6 class="card-title text-dark">${dapp.title}</h6>
                                <p class="card-text">${dapp.desc}</p>
                            </div>
                            <div class="card-footer text-right">
                                <button class="btn btn-sm btn-secondary text-uppercase">${$.i18n.prop('dapp_token_stay_tuned')}</button>
                            </div>
                        </div>
                    </div>
                `);
            }
        }

        var biz = {
            operation: "list",
        }
        Common.post('dapp/set', biz, {}, function (res) {
            if (res.base.code === "SUCCESS") {
                var dapps = res.biz;
                for (var i = 0; i < dapps.length; i++) {
                    var dapp = dapps[i];
                    $('.dapp-data').append(`
                    <div class="col-lg-3 col-md-4 col-sm-6" style="padding-bottom:15px;">
                        <div class="card shadow" >
                            <img src="${dapp.img}" with="390" style="max-height:280px"  class="card-img-top">
                            <div class="card-body" style="height:130px;">
                                <h6 class="card-title text-dark">${dapp.title}</h6>
                                <p class="card-text">${dapp.desc}</p>
                            </div>
                            <div class="card-footer text-right">
                                <button class="btn btn-danger btn-sm" onclick="DApps.removeDapp('`+dapp.id+`')">${$.i18n.prop('dapps_button_remove')}</button>
                                <a href="#" class="btn btn-sm btn-primary dapp-btn text-uppercase" dapp-name="${dapp.title}" dapp-url="${dapp.url}">${$.i18n.prop('dapps_button_enter')}</a>
                            </div>
                        </div>
                    </div>
                `);
                }
            }
        })

        $('.dapp-btn').bind('click', function () {
            var dappName = $(this).attr('dapp-name');
            var dappUrl = $(this).attr('dapp-url');
            var bodyp = $('.modal-body p').text();
            $('#myModal p').text(bodyp.replace(/GGGGG/g, dappName));

            $('.dapp-name').text(dappName);
            $('#myModal').modal('show');
            $('.modal-footer button:eq(1)').unbind('click').bind('click', function () {
                that.goBrowser(dappUrl);
            });
        });

    }
}
