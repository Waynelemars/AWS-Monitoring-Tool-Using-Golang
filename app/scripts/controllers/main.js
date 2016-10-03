'use strict';
/**
 * @ngdoc function
 * @name sbAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the sbAdminApp
 */
angular.module('sbAdminApp')
  .controller('MainCtrl', function($scope,$position, $http) {
	  
	$http.get('http://localhost:8080/metrics').success
	(function(response)
		{
		
		console.log(response);
		console.log(response.Reservations[0].Instances);
		$scope.loop = response.Reservations;
		$scope.reservations = response.Reservations[0].Instances;
		
		
		}
	)
	
  });
  
