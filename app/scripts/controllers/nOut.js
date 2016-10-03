'use strict';
/**
 * @ngdoc function
 * @name sbAdminApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the sbAdminApp
 
 */

angular.module('sbAdminApp')
  .controller('nOutCtrl', ['$scope', '$timeout', '$http' , function ($scope, $timeout,$http) {
	
	$http.get('http://localhost:8080/netwrkoutput').success
	(function(response)
		{
			var instanceArray = [];
			var dataArray = [];
			var metricArray = [];
			var pushArray = [];
			var dataArrayTwo =[];
			var metricArrayTwo =[];
			
			
			var count = Object.keys(response).length;
			console.log(count);
			
			for(var i=0 ; i< count ; i++) {
			
				console.log("Boom boom" , response[i]);
				instanceArray.push(response[i].InstanceId);
				if((response[i].Datapoints).length > 0) {
					
				var cnt = Object.keys(response[i].Datapoints).length; 
				console.log("DP cnt is", cnt);
				
			     }
				 
				if(i+1!=cnt) {
					for(var j=0 ; j < cnt;j++) {
					console.log("i is", i);
					console.log("j is", j);
					
						dataArray.push(response[i].Datapoints[j].Timestamp);
					    metricArray.push(response[i].Datapoints[j].Average);
			
				   }
				} 
				 else {
					 for(var j=0 ; j < cnt;j++) {
					 console.log("i is", i);
					 console.log("j is", j);
						//dataArray.push(response[i].Datapoints[j].Timestamp);
					    metricArrayTwo.push(response[i].Datapoints[j].Average);
						
				     }
				  }
			  }
			
			console.log("Instance array is", instanceArray);
			console.log("Data array ZZZ is", dataArray);
			console.log("metricArray is", metricArray);
			
			
		    $scope.bar = {
				labels: dataArray,
				series: instanceArray,

				data: [metricArray,metricArrayTwo ]
				
			}; 

		}
	)
	
   
}]);