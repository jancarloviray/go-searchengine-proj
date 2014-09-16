window.app = angular.module('app',['ngSanitize'])

app.controller('home', [
	'$scope', '$http','$sce',
	function($scope, $http, $sce){
		var model = {
			query: '',
			results: [],
			duration: '',
			values: []
		};
		$scope.$sce = $sce;
		$scope.query = "";
		$scope.model = model;
		$scope.$watch('query', function(n, o){
			if (n === "" || !n) {
				model.values = [];
				return
			};
			$http.get('/api/search?s=' + n, { cache: true, responseType: 'json'})
			.success(function(data, status){
				for(var i = 0, len = model.values.length; i < len; i++){
					if (!data.values[i]) {
						model.values = [];
						return
					};
					data.values[i] = data.values[i].replace(n, '<span class="hl">' + n + '</span>');
				}
				model.query = n;
				model.results = data.results;
				model.duration = data.duration;
				model.values = data.values;
			})
			.error(function(data, status){
			});
		});
	}]
);
