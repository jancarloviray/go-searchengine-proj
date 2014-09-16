window.app = angular.module('app',[])

app.controller('home', [
'$scope', '$http',
function($scope, $http){
	var model = {
		query: '',
		results: [],
		timing: '',
		values: []
	};
	$scope.query = "";
	$scope.model = model;
	$scope.$watch('query', function(n, o){
		$http.get('/api/search?s=' + n, { cache: true, responseType: 'json'})
		.success(function(data, status){
			model.query = n;
			model.results = data.results;
			model.timing = data.timing;
			model.values = data.values;
		})
		.error(function(data, status){
		});
	});
}]
);
