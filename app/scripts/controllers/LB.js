'use strict';
/**
 * @ngdoc function
 * @name sbAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the sbAdminApp
 */
angular.module('sbAdminApp')
  .controller('lBCtrl', function($scope,$position, $http) {
	  
	$http.get('http://localhost:8080/count').success
	(function(response)
		{
		$scope.LB = "LoadBalancer"
		$scope.count = response[0].Datapoints[0].Sum;
		
		
		}
	)
	
  });
  
