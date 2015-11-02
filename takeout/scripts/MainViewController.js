var API_URL = "http://localhost:8080";
var RESTAURANT_ID = '15c0e2b6-e73b-441d-8389-8f8942a60f08'

$(function () {
  // setup the data manager
  var _menuItemsDataManager = new MenuItemsDataManager();

  // setup knockout and the view controller
  var _mainViewController = new MainViewController(_menuItemsDataManager);
  ko.applyBindings(_mainViewController);
  _mainViewController.onLoad();

});


// This is a simple *viewmodel* - JavaScript that defines the data and behavior of your UI
function MainViewController(menuItemsDataManager) {
  var self = this;

  self.menuItems = ko.observableArray([]);

  self.newMenuItem = ko.observable(null);

  self.onLoad = function() {
    menuItemsDataManager.getAllMenuItems(RESTAURANT_ID, {
      onSuccess: function(menuItems) {
        self.menuItems(menuItems);
      }
    }, listen = true);
  }


  self.addMenuItem = function() {
    self.newMenuItem(menuItemsDataManager.getEmptyMenuItem());
  }

  self.saveMenuItem = function() {

    var item = self.newMenuItem();

    //TODO: ask the data manager for validation on menu item object

    menuItemsDataManager.createMenuItem(item, {
      onSuccess: function() {
        //the menuItems list should update itself automatically
      }
    })

  }

  self.cancelMenuItem = function () {
    self.newMenuItem(null);
  }


}
