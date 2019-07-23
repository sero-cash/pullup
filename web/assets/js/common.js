var Common = {
    host: 'http://127.0.0.1:2345',

    seroRpcHost:'',

    app: {},

    LANGUAGE_CODE: 'zh_CN',

    baseDecimal: new BigNumber(10).pow(18),

    init: function () {
        var that = this;
        that.app.init();
        that.getLang();

        $('.language').bind('click',function () {
            var lang_code = $.cookie('language');

            if ('zh_CN'===lang_code) {
                $.cookie('language', 'en_US');
                $('.language').text('简体中文');
            }else{
                $.cookie('language', 'zh_CN');
                $('.language').text('English');
            }

            that.getLang();
        });

    },

    getLang: function () {
        var _LANGUAGE_CODE
        if (!jQuery.i18n.normaliseLanguageCode({})){
            _LANGUAGE_CODE = jQuery.i18n.normaliseLanguageCode({}); //获取浏览器的语言
        }
        var lang_code = $.cookie('language');
        if (typeof lang_code === 'undefined') {
            lang_code = _LANGUAGE_CODE;
        }
        Common.app.loadProperties(lang_code);
    },

    post: function (_method, _biz, _page, callback) {
        var that = this;

        var result = new Object();
        var timestamp = 1234567;
        var sign = "67ff54447b89f06fe4408b89902e585167abad291ec41118167017925e24e320";
        var data = {
            base: {
                timestamp: timestamp,
                sign: sign,
            },
            biz: _biz,
            page: _page,
        }

        $.ajax({
            url: that.host + '/' + _method,
            type: 'post',
            dataType: 'json',
            async: false,
            data: JSON.stringify(data),
            beforeSend: function () {
            },
            success: function (res) {
                if (callback) {
                    callback(res)
                }
            }
        })

        return result;
    },

    postAsync: function (_method, _biz, _page, callback) {
        var that = this;

        var result = new Object();
        var timestamp = 1234567;
        var sign = "67ff54447b89f06fe4408b89902e585167abad291ec41118167017925e24e320";
        var data = {
            base: {
                timestamp: timestamp,
                sign: sign,
            },
            biz: _biz,
            page: _page,
        }

        $.ajax({
            url: that.host + '/' + _method,
            type: 'post',
            dataType: 'json',
            async: true,
            data: JSON.stringify(data),
            beforeSend: function () {
            },
            success: function (res) {
                if (callback) {
                    callback(res)
                }
            }
        })

        return result;
    },


}

