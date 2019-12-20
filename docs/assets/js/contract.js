var ContractDashboard = {

    init: function () {
        var that = this;

        that.showAll();

    },

    showAll: function () {


    },

    watchContract: function () {

    },

    watchToken: function () {

        var innerHtml = `
            <form>
              <div class="form-group">
                <label for="exampleFormControlSelect1">FROM</label>
                <select class="form-control" id="exampleFormControlSelect1">
                  <option>1</option>
                  <option>2</option>
                  <option>3</option>
                  <option>4</option>
                  <option>5</option>
                </select>
              </div>
              <div class="form-group">
                <label for="exampleFormControlTextarea1">Token Name</label>
                <input class="form-control" id="name" />
              </div>
              <div class="form-group">
                <label for="exampleFormControlTextarea1">Token Symbol</label>
                <input class="form-control" id="symbol" />
              </div>
              <div class="form-group">
                <label for="exampleFormControlTextarea1">Token Decimals</label>
                <input class="form-control" id="decimal" />
              </div>
            </form>
        `;

        $('.modal-title').empty().text("Add Token");
        $('.modal-body').empty().append(innerHtml);

    }

}

var ContractDeploy = {

    init: function () {

    },

}

var ContractDetail = {

    init: function () {

    },

}