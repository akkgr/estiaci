"use strict";angular.module("newFolderApp",["ngAnimate","ngCookies","ngResource","ngRoute","ngSanitize","ngTouch","angular-loading-bar"]).config(["$routeProvider",function(a){a.when("/",{templateUrl:"views/main.html",controller:"MainCtrl",controllerAs:"main"}).when("/buildings",{templateUrl:"views/buildings.html",controller:"BuildingsCtrl",controllerAs:"buildings"}).when("/buildings/:id",{templateUrl:"views/building.html",controller:"BuildingCtrl",controllerAs:"building"}).otherwise({redirectTo:"/"})}]).directive("lastRepeaterElement",function(){return function(a){a.$last&&a.$emit("LastRepeaterElement")}}).controller("NavigationCtrl",["$scope","$location",function(a,b){a.isCurrentPath=function(a){return b.path()==a}}]),angular.module("newFolderApp").controller("MainCtrl",function(){this.awesomeThings=["HTML5 Boilerplate","AngularJS","Karma"]}),angular.module("newFolderApp").controller("BuildingsCtrl",["$scope","$http",function(a,b){a.data=[],b({method:"GET",url:"/api/buildings"}).then(function(b){a.data=b.data},function(a){toastr.error("Σφάλμα: "+a.statusText)}),a.$on("LastRepeaterElement",function(){$('[data-toggle="tooltip"]').tooltip()})}]),angular.module("newFolderApp").controller("BuildingCtrl",["$scope","$http","$routeParams",function(a,b,c){var d=c.id,e=new google.maps.LatLng(37.9667,23.7167),f={zoom:15,center:e,mapTypeId:google.maps.MapTypeId.ROADMAP},g=new google.maps.Map(document.getElementById("map"),f),h=(new google.maps.Geocoder,new google.maps.Marker({map:null,position:e,title:""}));0!=d&&b({method:"GET",url:"/api/buildings/"+a.id}).then(function(b){a.data=b.data;var c=new google.maps.LatLng(a.data.address.location.Lat,a.data.address.location.Lng);g.setCenter(c),h.setPosition(c),h.setTitle(adr)},function(a){toastr.error("Σφάλμα: "+a.statusText)}),a.change=function(b){var c=a.data.address.street+" "+a.data.address.streetnumber+", "+a.data.address.postalcode+" "+a.data.address.area;a.geocoder.geocode({address:c},function(b,d){if(d===google.maps.GeocoderStatus.OK){a.data.address.location={Lng:b[0].geometry.location.lng(),Lat:b[0].geometry.location.lat()};var f=new google.maps.LatLng(a.data.address.location.Lat,a.data.address.location.Lng);g.setCenter(f),h.setPosition(f),h.setTitle(c)}else g.setCenter(e),h.setMap(null),h.setTitle(""),toastr.error("Geocode was not successful for the following reason: "+d)}),null!=b&&b()},a.save=function(){if(0!=d)var c={method:"POST",url:"/api/buildings",data:JSON.stringify(a.data)};else var c={method:"PUT",url:"/api/buildings/"+d,data:JSON.stringify(a.data)};b(c).then(function(b){a.data=b.data,d=b.data.Id,toastr.success("Έγινε αποθήκευση.");var c=a.data.address.location,e=new google.maps.LatLng(c.Lat,c.Lng);g.setCenter(e),h.setPosition(e),h.setTitle(adr)},function(a){toastr.error("ΔΕΝ έγινε αποθήκευση. Σφάλμα: "+a.statusText)})}}]),angular.module("newFolderApp").run(["$templateCache",function(a){a.put("views/building.html",'<div class="row-fluid"> <div class="col-xs-6 col-md-6"> <form class="form-horizontal"> <fieldset> <div class="panel panel-default"> <div class="panel-heading">Διεύθυνση</div> <div class="panel-body"> <div class="form-group"> <label class="col-md-4 control-label" for="street">οδός</label> <div class="col-md-4"> <input id="street" name="street" type="text" placeholder="οδός" class="form-control input-md" ng-model="data.address.street"></div> <div class="col-md-4 text-right"> <button id="mapbutton" name="mapbutton" class="btn btn-primary" ng-click="change()">Χάρτης</button> </div> </div> <div class="form-group"> <label class="col-md-4 control-label" for="streetnumber">αριθμός</label> <div class="col-md-2"> <input id="streetnumber" name="streetnumber" type="text" placeholder="αριθμός" class="form-control input-md" ng-model="data.address.streetnumber"></div> </div> <div class="form-group"> <label class="col-md-4 control-label" for="area">περιοχή</label> <div class="col-md-4"> <input id="area" name="area" type="text" placeholder="περιοχή" class="form-control input-md" ng-model="data.address.area"></div> </div> <div class="form-group"> <label class="col-md-4 control-label" for="postalcode">Τ.Κ.</label> <div class="col-md-2"> <input id="postalcode" name="postalcode" type="text" placeholder="Τ.Κ." class="form-control input-md" ng-model="data.address.postalcode"></div> </div> </div> </div> <div class="panel panel-default"> <div class="panel-heading">Πληροφορίες</div> <div class="panel-body"> <div class="form-group"> <div class="checkbox col-md-2 col-md-offset-2"> <input class="styled" type="checkbox" ng-model="data.managment"> <label>Διαχείρηση</label> </div> <div class="checkbox col-md-2"> <input class="styled" type="checkbox" ng-model="data.active"> <label>Ενεργή</label> </div> <label class="col-md-4 control-label">Διαμερίσματα</label> <div class="col-md-2"> <button id="mapbutton" name="mapbutton" class="btn btn-default">{{data.appartments.length}}</button> </div> </div> </div> </div> <div class="col-md-4 col-md-offset-8 text-right"> <button id="savebutton" name="savebutton" class="btn btn-success" ng-click="save()">Αποθήκευση</button> </div> </fieldset> </form> </div> <div class="col-xs-6 col-md-6" id="map" style="height: 520px"></div> </div>'),a.put("views/buildings.html",'<table class="table table-hover"> <thead> <tr> <th>Οδός</th> <th>Αριθμός</th> <th>Περιοχή</th> <th>T.K.</th> <th>Διαχείρηση</th> <th>Ενεργή</th> <th> <a ng-href="#/buildings/0" data-toggle="tooltip" data-placement="top" title="Προσθήκη νέου κτιρίου"> <span class="glyphicon glyphicon-plus-sign"></span> </a> </th> </tr> <tr> <td> <input type="text" ng-model="search.address.street"></td> <td> <input type="text" ng-model="search.address.streetnumber"></td> <td> <input type="text" ng-model="search.address.area"></td> <td> <input type="text" ng-model="search.address.postalcode"></td> <td> <div class="checkbox"> <input class="styled" type="checkbox" ng-model="search.managment"> <label></label> </div> </td> <td> <div class="checkbox"> <input class="styled" type="checkbox" ng-model="search.active"> <label></label> </div> </td> <td>{{(data | filter: search).length}} από {{data.length}}</td> </tr> </thead> <tbody> <tr ng-repeat="item in data | filter: search | orderBy: [\'address.street\',\'address.streetnumber\',\'address.area\'] as filtered_result track by item.id" last-repeater-element> <td>{{item.address.street}}</td> <td>{{item.address.streetnumber}}</td> <td>{{item.address.area}}</td> <td>{{item.address.postalcode}}</td> <td> <div class="checkbox"> <input class="styled" type="checkbox" ng-model="item.managment" disabled> <label></label> </div> </td> <td> <div class="checkbox"> <input class="styled" type="checkbox" ng-model="item.active" disabled> <label></label> </div> </td> <td> <a ng-href="#/buildings/{{item.id}}" data-toggle="tooltip" data-placement="top" title="Καρτέλα κτιρίου"> <span class="glyphicon glyphicon-edit"></span> </a> </td> </tr> </tbody> </table>'),a.put("views/main.html",'<div class="jumbotron"> <h1>\'Allo, \'Allo!</h1> <p class="lead"> <img src="images/yeoman.c582c4d1.png" alt="I\'m Yeoman"><br> Always a pleasure scaffolding your apps. </p> <p><a class="btn btn-lg btn-success" ng-href="#/">Splendid!<span class="glyphicon glyphicon-ok"></span></a></p> </div> <div class="row-fluid marketing"> <h4>HTML5 Boilerplate</h4> <p> HTML5 Boilerplate is a professional front-end template for building fast, robust, and adaptable web apps or sites. </p> <h4>Angular</h4> <p> AngularJS is a toolset for building the framework most suited to your application development. </p> <h4>Karma</h4> <p>Spectacular Test Runner for JavaScript.</p> </div>')}]);