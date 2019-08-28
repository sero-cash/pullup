var DApps = {

    init: function () {
        var that = this;
        setTimeout(function () {
            that.genPageData();
        }, 10)

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

                $('.modal-title').text($.i18n.prop('dapps_modal_title'));
                $('.modal-body p').text($.i18n.prop('dapps_modal_body'));
                $('.modal-footer button:eq(0)').text($.i18n.prop('dapps_button_cancel'));
                $('.modal-footer button:eq(1)').text($.i18n.prop('dapps_button_enter'));

            }
        });

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
            url: "http://134.175.161.78:8088",
            showTips: true,
            state: 1,
        },
        {
            img: "./assets/img/sanguo.png",
            title: "超零三国-无限穿越",
            desc: "",
            author: "盘古",
            url: "",
            showTips: true,
            state: 0,
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
            url: "http://134.175.161.78:8088",
            showTips: true,
            state: 1,
        },
        {
            img: "./assets/img/sanguo.png",
            title: "超零三国-无限穿越",
            desc: "",
            author: "盘古",
            url: "",
            showTips: true,
            state: 0,
        }
    ],

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
            if(dapp.state === 1){
                $('.dapp-data').append(`
                    <div class="col-lg-3">
                        <div class="card shadow">
                            <img src="${dapp.img}" class="card-img-top">
                            <div class="card-body" style="height:200px;">
                                <h6 class="card-title text-dark">${dapp.title}</h6>
                                <p class="card-text">${dapp.desc}</p>
                            </div>
                            <div class="card-footer text-right">
                                <a href="${dapp.showTips ? "#" : dapp.url}" class="btn btn-sm btn-primary dapp-btn" dapp-name="${dapp.title}" dapp-url="${dapp.url}">${$.i18n.prop('dapps_button_enter')}</a>
                            </div>
                        </div>
                    </div>
                `);
            }else if(dapp.state === 0){
                $('.dapp-data').append(`
                    <div class="col-lg-3">
                        <div class="card shadow">
                            <img src="${dapp.img}" class="card-img-top">
                            <div class="card-body" style="height:200px;">
                                <h6 class="card-title text-dark">${dapp.title}</h6>
                                <p class="card-text">${dapp.desc}</p>
                            </div>
                            <div class="card-footer text-right">
                                <button class="btn btn-sm btn-secondary">${$.i18n.prop('dapp_token_stay_tuned')}</button>
                            </div>
                        </div>
                    </div>
                `);
            }

        }

        $('.dapp-btn').bind('click', function () {
            var dappName = $(this).attr('dapp-name');
            var dappUrl = $(this).attr('dapp-url');
            var bodyp = $('.modal-body p').text();
            $('.modal-body p').text(bodyp.replace(/GGGGG/g, dappName));

            $('.dapp-name').text(dappName);
            $('.modal').modal('show');
            $('.modal-footer button:eq(1)').unbind('click').bind('click', function () {
                window.location.href = dappUrl;
            });
        });

    }
}