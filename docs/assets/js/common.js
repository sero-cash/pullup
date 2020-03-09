var Common = {
    host: 'http://127.0.0.1:2345',

    seroRpcHost: '',

    app: {},

    LANGUAGE_CODE: 'zh_CN',

    baseDecimal: new BigNumber(10).pow(18),

    init: function () {
        var that = this;
        setTimeout(function () {
            that.app.init();
            that.getLang();
            $('.language').bind('click', function () {
                var lang_code = $.cookie('language');

                if ('zh_CN' === lang_code) {
                    $.cookie('language', 'en_US');
                    $('.language').text('简体中文');
                } else {
                    $.cookie('language', 'zh_CN');
                    $('.language').text('English');
                }

                that.getLang();
            });
        }, 100);

        setTimeout(function () {
            that.checkVersion();
        },1000)
    },

    checkVersion :function(){
        var forceUpdateVersions = ["v0.1.13","v0.1.14"];
        var latestVersion = "v0.1.15";
        var that = this;
        $.ajax({
            url: that.host + '/version',
            type: 'get',
            dataType: 'json',
            async: false,
            success: function (serverVersion) {
                $('.version').text(serverVersion)
                if("EOF"===serverVersion || forceUpdateVersions.indexOf(serverVersion)>-1){
                    var localUtc = new Date().getTimezoneOffset() / 60;
                    if (localUtc === -8){
                        $('.update-body').empty().append(`
                            <ul class="list-group text-left">
                                <li class="list-group-item">SERO Pullup钱包已升级至${latestVersion}版本，之前的版本将停止使用。请下载最新的钱包：</li>
                                <li class="list-group-item">
                                    MacOS x64: <a href="https://sero-media-1256272584.cos.ap-shanghai.myqcloud.com/pullup/${latestVersion}/pullup-mac-${latestVersion}-zh_CN.tar.gz" target="_blank">pullup-mac-${latestVersion}-zh_CN.tar.gz</a>
                                </li>
                                <li class="list-group-item">
                                    Windows(PC): <a href="https://sero-media-1256272584.cos.ap-shanghai.myqcloud.com/pullup/${latestVersion}/pullup-windows-${latestVersion}-zh_CN.zip" target="_blank">pullup-windows-${latestVersion}.zip</a>
                                </li>
                            </ul>
                            `);
                    }else{
                        $('.update-body').empty().append(`
                            <ul class="list-group text-left">
                                <li class="list-group-item">The SERO Pullup wallet has been upgraded to version ${latestVersion} and the previous version will be discontinued. Please download the latest wallet:</li>
                                <li class="list-group-item">
                                    MacOS x64: <a href="https://github.com/sero-cash/pullup/releases/download/${latestVersion}/pullup-mac-${latestVersion}.tar.gz">pullup-mac-${latestVersion}.tar.gz</a></li>
                                <li class="list-group-item">
                                    Windows(PC): <a href="https://github.com/sero-cash/pullup/releases/download/${latestVersion}/pullup-windows-${latestVersion}.zip" target="_blank">pullup-windows-${latestVersion}.tar.gz</a>
                                </li>
                            </ul>
                        `);
                    }
                    $('#updateModal').modal({backdrop: 'static', keyboard: false});
                }else{
                    //check remote version.json
                    if (serverVersion!="v0.1.14"){
                        $.ajax({
                            url: that.host + '/remoteVersion',
                            type: 'get',
                            dataType: 'json',
                            async: false,
                            success: function (remoteVersion) {
                                if(serverVersion !== remoteVersion.version.app){
                                    var localUtc = new Date().getTimezoneOffset() / 60;
                                    var title = '';
                                    var desc = [];
                                    if (localUtc === -8){
                                        title = `版本${remoteVersion.version.app}更新内容：`;
                                        desc = remoteVersion.description.zh;

                                    }else{
                                        title = `${remoteVersion.version.app} updated features：`;
                                        desc = `${remoteVersion.description.en}`
                                    }
                                    var content = '';
                                    for(var i=0;i<desc.length;i++){
                                        content += `<p class="text-info">${desc[i]}</p>`
                                    }

                                    $('.update-body').empty().append(`
                                            <ul class="list-group text-left">
                                                <li class="list-group-item">${title}</li>
                                                <li class="list-group-item">
                                                    ${content}
                                                </li>
                                                <li class="list-group-item">
                                                    MacOS x64: <a href="${remoteVersion.version.appUrl.mac}" target="_blank">pullup-mac-${remoteVersion.version.app}-zh_CN.tar.gz</a>
                                                </li>
                                                <li class="list-group-item">
                                                    Windows(PC): <a href="${remoteVersion.version.appUrl.win}" target="_blank">pullup-windows-${remoteVersion.version.app}.zip</a>
                                                </li>
                                            </ul>
                                        `);
                                    $('#updateModal').modal({backdrop: 'static', keyboard: false});
                                }
                            }
                        })
                    }
                }
            }
        })
    },

    convertErrors: function (err) {
        var s = err;
        if (err.indexOf("could not decrypt key with given passphrase") > -1) {
            s = $.i18n.prop('convert_error_password');
        } else if (err.indexOf("no enough unlocked utxos") > -1) {
            s = $.i18n.prop('convert_error_utxo');
        } else if (err.indexOf("stx Verify error") > -1) {
            s = $.i18n.prop('convert_error_verifytx');
        }
        return s
    },

    getLang: function () {
        var _LANGUAGE_CODE = "en_US";
        // if (!jQuery.i18n.normaliseLanguageCode({})) {
        //     _LANGUAGE_CODE = jQuery.i18n.normaliseLanguageCode({}); //获取浏览器的语言
        // }
        var lang_code = $.cookie('language');
        if (!lang_code) {
            lang_code = _LANGUAGE_CODE;
        }
        if ('zh_CN' === lang_code) {
            $('.language').text('English');
        } else {
            $('.language').text('简体中文');
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

    //_params is an array
    postSeroRpc: function (_method, _params, callback) {
        var that = this;
        var postData = {
            id: 0,
            jsonrpc: "2.0",
            method: _method,
            params: _params,
        };

        $.ajax({
            url: that.host + '/rpc',
            type: 'post',
            dataType: 'json',
            async: true,
            data: JSON.stringify(postData),
            beforeSend: function () {
            },
            success: function (res) {
                if (callback) {
                    callback(res)
                }
            }
        })
    },

    //_params is an array
    postSeroRpcSync: function (_method, _params, callback) {
        var that = this;
        var postData = {
            id: 0,
            jsonrpc: "2.0",
            method: _method,
            params: _params,
        };

        $.ajax({
            url: that.host + '/rpc',
            type: 'post',
            dataType: 'json',
            async: false,
            data: JSON.stringify(postData),
            beforeSend: function () {
            },
            success: function (res) {
                if (callback) {
                    callback(res)
                }
            }
        })
    },

    //_params is an array
    postPullupRpc: function (_method, _params, callback) {
        var that = this;
        var postData = {
            id: 0,
            method: _method,
            params: _params,
        };

        $.ajax({
            url: that.host + '/pullup_rpc',
            type: 'post',
            dataType: 'json',
            async: true,
            data: JSON.stringify(postData),
            beforeSend: function () {
            },
            success: function (res) {
                if (callback) {
                    callback(res)
                }
            }
        })
    },

    


}

