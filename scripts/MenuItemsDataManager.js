function MenuItemsDataManager() {
  var self = this;

  self.getEmptyMenuItem = function() {
    return {
      MenuItemID: null,
      Name: null,
      Description:null,
      Price:null,
      ImageURL:null,
      RestaurantID: RESTAURANT_ID
    }
  }

  var _cachedMenuItems = null;
  self.getAllMenuItems = function(restaurantID, callbacks) {
    if(_cachedMenuItems != null) {
      callbacks.onSuccess(_cachedMenuItems);
      return;
    }

    $.ajax({
      method: "GET",
      url: API_URL + "/menu/items?restaurantID=" + restaurantID,
      contentType:"application/json"
    })
    .success(function(data,textStatus, jqXHR) {
      console.log(data);
      _cachedMenuItems = data;
      callbacks.onSuccess(_cachedMenuItems);
    })
    .error(function(data,textStatus, jqXHR) {
      console.log("ERROR: getting all menu items for restaurantID " + restaurantID);
      if(callbacks.onError) {
        callbacks.onError(data);
      }
    });
  }


  self.createMenuItem = function(item, callbacks) {

    console.log(item);

    $.ajax({
      method: "POST",
      url: API_URL + "/menu/items",
      contentType:"application/json",
      data: JSON.stringify(item)
    })
    .success(function(data,textStatus, jqXHR) {

      //will need to add an ID to the item here
      _cachedMenuItems.push(item);
      callbacks.onSuccess(item);
    })
    .error(function(data,textStatus, jqXHR) {
      console.log("ERROR: creating new menu item");
      if(callbacks.onError) {
        callbacks.onError(data);
      }
    });

  }

}
