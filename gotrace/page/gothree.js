/**
 * GoThree.js - generate three.js code for Go concurrency data.
 */

var GoThree = GoThree || (function () {
	// clear() resets the scene;
	THREE.Object3D.prototype.clear = function(){
		var children = this.children;
		for(var i = children.length-1;i>=0;i--){
			var child = children[i];
			this.remove(child);
		};
	};

	return {}
})();

GoThree.Trace = function() {
	var _data = {};
	var _scene = undefined;

	var _colors1 = {
		"send_arrow":  "#BDF271",
		"send_value":  "#BDF271",
		"go_normal":  "#29D9C2",
		"go_blocked":  "#01A2A6",
		"go_sleep":  "#01A2A6",
		"go_link":  "#29D9C2",
		"go_cap":  "#615469"
	};
	var _colors2 = {
        "send_arrow":  "#01B0F0",
        "send_value":  "#01B0F0",
        "go_normal":  "#AEEE00",
        "go_link":  "#AEEE00",
        "go_blocked":  "#FF358B",
        "go_sleep":  "#FF358B",
        "go_cap":  "#424242",
    };
	var _colors3 = {
        "send_arrow":  "#00D6DD",
        "send_value":  "#00D6DD",
        "go_normal":  "#D4FF00",
        "go_link":  "#B1B4B5",
        "go_blocked":  "#ED0000",
        "go_sleep":  "#6C8200",
        "go_cap":  "#D4FF00",
    };

	var _colors = _colors3;

	// Predefined materials
	var _mats = {
		"go_normal": new THREE.LineBasicMaterial( { color: _colors["go_normal"], linewidth: 5 } ),
		"go_sleep": new THREE.LineBasicMaterial( { color: _colors["go_sleep"], linewidth: 2 } ),
		"go_blocked": new THREE.LineBasicMaterial( { color: _colors["go_blocked"], linewidth: 2 } ),
		"go_link": new THREE.LineBasicMaterial( { color: _colors["go_link"], linewidth: 1, } ),
		"go_cap": new THREE.MeshBasicMaterial({color: _colors["go_cap"]}),
		"send_value": new THREE.MeshBasicMaterial({color: _colors["send_value"]}),
		"send_arrow": new THREE.LineBasicMaterial({color: _colors["send_arrow"]}),
		"channel": new THREE.MeshBasicMaterial( {color: _colors["channel"]} ),
	};


	var _t = 0;
	var _scale = 1;

	// total drawing time/length. Actually, main is of length _total_time.
	var _total_time = 100;

	// initY is a start Y point (where main starts)
	var initY = _total_time / 2;

	var _run_time = 0;
	var _tick_per_sec = 0;
	var _speed = 1; // 1/x time multiplier
	var _speedup_factor = 1.5;

	// _step represents current step in iterating data. 0 <= x < data.length
	var _step = 0;

	// _goroutines represents info&objects for active goroutines
	var _goroutines = [];

	// _gocaps holds text objects for Gorotuines caps.
	var _gocaps = [];

	// _traces represents info&objects for trances of inter-channel sends
	var _traces = [];

	// _cam_position holds current camera/controls positions.
	var _cam_position = {};

	// misc params from the caller.
	var _params = {};

	// distance between goroutines.
	var _distance = 100;

	// distanceShift on each goroutine creation.
	var _distanceShift = 0;

	// global angleShift.
	var _angleShift = 0;
	var _autoAngle = false;

	// Reset timing, but keeps data. Aims for restarting animation.
	this.resetTime = function() {
		_scene.clear();

		_t = 0;
		_step = 0;
		_goroutines = [];
		_traces = [];
		_angleShift = 0;
		_distanceShift = 0;
	};

	// Init w/ scene and new data
	this.init = function(scene, data, params, scale) {
		_scale = scale;
		_scene = scene;
		_data = data;
		_params = params;
		if (_params.speed !== undefined) {
			_speed = _params.speed;
		}
		if (_params.distance !== undefined) {
			_distance = _params.distance;
		}
		if (_params.autoAngle !== undefined) {
			_autoAngle = _params.autoAngle;
		}
		if (_params.totalTime !== undefined) {
			_total_time = _params.totalTime;
		}
		this.resetTime();
		this._updateTimings();
	};

	// Update timing related calculations
	this._updateTimings = function() {
		_run_time = _data[_data.length-1].t;
		// FIXME: msecs_per_tick is more proper name
		_tick_per_sec = _total_time/_run_time;
	};

	this.animate = function() {
		TWEEN.update();
		if (_data.length === 0) return;
		if (_t > _total_time || _step > _data.length) return;
		
		// Process next command
		var cmd = _data[_step];
		if (cmd === undefined) return;
		// if current command is ready to be rendered...
		while (_t/_tick_per_sec > cmd.t) {
			// process command
			switch (cmd.command) {
				case 'create goroutine':
					var pos = this._calculatePosition(cmd.name, cmd.parent, -_t);
					if (pos === undefined) {
						pos = {x: 0, y: 0, z: 0};
					}
					this._cmd_create_goroutine(cmd.name, cmd.parent, pos);
				break;
				case 'create channel':
					this._cmd_create_channel(cmd.name, cmd.goroutine);
				break;
				case 'send to channel':
					this._cmd_send_to_channel(cmd.ch, cmd.from, cmd.to, cmd.value, cmd.duration);
				break;
				case 'stop goroutine':
					this._cmd_stop_goroutine(cmd.name);
				break;
				case 'block goroutine':
					this._cmd_block_goroutine(cmd.name);
				break;
				case 'unblock goroutine':
					this._cmd_unblock_goroutine(cmd.name);
				break;
				case 'sleep goroutine':
					this._cmd_sleep_goroutine(cmd.name);
				break;
			}
			_step++
			cmd = _data[_step];
			if (cmd === undefined) break;
		}

		// increase time
		_t += 1/_speed;

		// grow existing goroutines along the time axis
		for (var i = 0; i < _goroutines.length; i++) {
			var geom = _goroutines[i].line.geometry;
			var end = geom.vertices[1];
			end.y -= _scale*(1/_speed);
			geom.verticesNeedUpdate = true;
		};
	};

	// 'create goroutine' command handler
	this._cmd_create_goroutine = function(name, parent, pos) {
		// create new goroutine object
		var goroutine = {
			name: name,
			parent: parent,
		};

		var x = pos.x;
		var y = initY + pos.y;
		var z = pos.z;

		var depth = this._getDepth(parent);
		var width = 5;
		if (depth < 5) {
			width = width-depth;
		}

		// create new line
		var geom = new THREE.Geometry();
		var start = new THREE.Vector3( x, y, z );
		var end = new THREE.Vector3( x, y, z );
		geom.vertices.push( start, end );
		goroutine.line = new THREE.Line(geom, _mats["go_normal"]);

		// create link with parent line
		if (parent != undefined) {
			var p = _goroutines.find({name: parent});
			if (p !== undefined) {
				var pgeom = p.line.geometry.vertices[0];
				var lstart = new THREE.Vector3(pgeom.x, y, pgeom.z);
				var lend = new THREE.Vector3( x, y, z );
				var lgeom = new THREE.Geometry();
				lgeom.vertices.push( lstart, lend );
				lline = new THREE.Line( lgeom, _mats["go_link"] );
				_scene.add( lline );
			}
		}

		// create cap text
		var shapes = THREE.FontUtils.generateShapes(this._goroutine_name(name), { font: "helvetiker", weight: "normal", size: 3*_scale } );
		var geom = new THREE.ShapeGeometry( shapes );
		var text = new THREE.Mesh( geom, _mats["go_cap"] );
		text.position.x = x;
		text.position.y = y;
		text.position.z = z;
		text.lookAt( _cam_position );

		// add to scene
		_scene.add( text );
		_scene.add( goroutine.line );

		_gocaps.push(text);
		_goroutines.push(goroutine);
	};

	// 'create goroutine' command handler
	this._cmd_send_to_channel = function(ch, from, to, value, dur) {
		var s = _goroutines.find({name: from});
		var e = _goroutines.find({name: to});
		if (s === undefined) {
			console.log("ERR: end of line is undefined for "+from);
			return;
		}
		if (e === undefined) {
			console.log("ERR: end of line is undefined for "+to);
			return;
		}
		var start = s.line.geometry.vertices[1];
		var end = e.line.geometry.vertices[1];

		var duration = 50; // default: 0.05s
		if (dur !== undefined) {
			duration = dur / (1000 * 1000); // -> nanoseconds -> milliseconds
		}

		// target end is on the same Y axis as the start
		var target = { x: end.x, y: end.y, z: end.z };
		var targetV = new THREE.Vector3(end.x, end.y, end.z);

		// First, create line that will be animated
		// (as ArrowHelper can't be animated)
		var head = { x: start.x, y: end.y, z: start.z };
		var tail = { x: start.x, y: end.y, z: start.z };
		var geom = new THREE.Geometry();
		geom.vertices.push(head);
		geom.vertices.push(tail);
		var line = new THREE.Line(geom, _mats["send_arrow"]);
		_scene.add(line); 

		// create tween for line animation
		var d = duration/1000000*_tick_per_sec;
		var d2 = Math.sqrt(d*d + _distance*_distance);
		var tween = new TWEEN.Tween(tail).to(target, d2 * _speed * 2);
		tween.easing(TWEEN.Easing.Cubic.InOut);

		// create text with value
		var trace = {};
		var shapes = THREE.FontUtils.generateShapes(value, { font: "helvetiker", weight: "normal", size: 1.5*_scale } );
		var geom = new THREE.ShapeGeometry( shapes );
		trace.text = new THREE.Mesh( geom, _mats["send_value"] );
		// get center of the line
		var c = new THREE.Vector3().addVectors( start, targetV ).multiplyScalar( 0.5 );
		trace.text.position.x = c.x;
		trace.text.position.y = c.y;
		trace.text.position.z = c.z;
		trace.text.lookAt( _cam_position );
		_scene.add(trace.text);

		_traces.push(trace);

		// Tweening
		tween.onUpdate(function(){
			line.geometry.vertices[1].x = tail.x;
			line.geometry.vertices[1].y = tail.y;
			line.geometry.vertices[1].z = tail.z;
			line.geometry.verticesNeedUpdate = true;
		}).onComplete(function(){
			var direction = new THREE.Vector3().subVectors(targetV, head);
			var len = direction.length();
			var arrow = new THREE.ArrowHelper(direction.clone().normalize(), head, len, _colors["send_arrow"], 1, 1);
			_scene.add( arrow ); 
		});

		tween.start();
	};

	// _goroutine_name formats goroutine name.
	this._goroutine_name = function(name) {
		if (name == "#1") {
			return "main";
		}

		// is asked, draw caps text for all goroutines
		if (_params.allCaps == true) {
			return name;
		}
		return "";
	};

	// onControlsChanged runs when controls/camera position is changed.
	this.onControlsChanged = function(controls) {
		_cam_position = controls.position;
		this._updateTextFaces();
	};

	// _updateTextFaces rotates texts to look at camera
	this._updateTextFaces = function() {
		_gocaps.each(function(g) {
			g.lookAt( _cam_position );
		});

		_traces.each(function(g) {
			g.text.lookAt( _cam_position );
		});
	};

	// slowdown slows down animation.
	this.slowdown = function() {
		_speed = _speed * _speedup_factor;;
	};

	// speedup speeds up animation.
	this.speedup = function() {
		_speed = _speed / _speedup_factor; 
	};

	this._getDepth = function(parent) {
		var depth = 0;
		var node = _goroutines.find({name: parent});
		while (node != undefined) {
			depth++;
			node = _goroutines.find({name: node.parent});
		}
		return depth
	}


	this._calculatePosition = function(name, parent, y) {
		var siblings = _goroutines.findAll({parent: parent});

		// calculate goroutine depth
		var depth = this._getDepth(parent);

		var position = {x:0, y:y, z:0};
		if (parent != undefined) {
			var p = _goroutines.find({name: parent});
			if (p === undefined) 
				return p;
			var parentStart = p.line.geometry.vertices[0];

			var distance = _distance/depth;
			if (depth > 1 && _params.distance2 != undefined) {
				distance = _params.distance2;
			}
			if (_params.autoGrow == true) {
				distance += _distanceShift;
				_distanceShift += 0.2;
			}

			// scale
			distance *= _scale;

			// calculate parent's angle
			var initAngle = _angleShift;
			var grandParent = _goroutines.find({name: p.parent});
			if (grandParent !== undefined) {
				var p1 = parentStart;
				var p2 = grandParent.line.geometry.vertices[0];
				var vec = new THREE.Vector3(p1.x-p2.x, 0, p1.z-p2.z);
				var orig = new THREE.Vector3(1, 0, 0);
				initAngle = initAngle + orig.angleTo(vec) * 180/Math.PI;
			}

			var singleAngle = 15;
			if (_params.angle !== undefined && depth == 1) {
				singleAngle = _params.angle;
			} else if (_params.angle2 !== undefined && depth == 2) {
				singleAngle = _params.angle2;
			}
			var angle = initAngle + singleAngle * siblings.length;
			if (_autoAngle) _angleShift += singleAngle;

			position.x = parentStart.x + Math.cos(angle * Math.PI/180)*distance;
			position.z = parentStart.z + Math.sin(angle * Math.PI/180)*distance;
		} 
		return position
	};

	// 'stop goroutine' command handler
	this._cmd_stop_goroutine = function(name) {
		var g = _goroutines.find({name: name});
		if (g === undefined) return;
		if (g.parent !== undefined) {
			var parent = _goroutines.find({name: g.parent});
			if (parent === undefined) {
				_goroutines.remove({name: name});
				return;
			}
			var ggeom = g.line.geometry.vertices[1];
			var pgeom = parent.line.geometry.vertices[1];
			var start = new THREE.Vector3(ggeom.x, ggeom.y, ggeom.z);
			var end = new THREE.Vector3(pgeom.x, ggeom.y, pgeom.z);
			var geom = new THREE.Geometry();
			geom.vertices.push( start, end );
			line = new THREE.Line( geom, _mats["go_link"] );
			_scene.add( line );
			_goroutines.remove({name: name});
		}
	};

	this._cmd_create_channel = function(name, goroutine) {
		var g = _goroutines.find({name: goroutine});
		if (g === undefined) {
			console.log("WARN: creating channel on unexistant goroutine: "+goroutine);
			return;
		}
		var pos = g.line.geometry.vertices[1];
		var geom = new THREE.SphereGeometry( 1.2, 1, 1 );
		var sphere = new THREE.Mesh( geom, _mats["channel"] );
		sphere.position.x = pos.x;
		sphere.position.y = pos.y;
		sphere.position.z = pos.z;
		_scene.add( sphere );
	};

	this._cmd_block_goroutine = function(name) {
		this._change_g_color(name, _mats["go_blocked"]);
	};

	this._cmd_unblock_goroutine = function(name) {
		this._change_g_color(name, _mats["go_normal"]);
	};

	this._cmd_sleep_goroutine = function(name) {
		this._change_g_color(name, _mats["go_sleep"]);
	};

	this._change_g_color = function(name, mat) {
		var g = _goroutines.find({name: name});
		if (g === undefined) return;

		var og = g.line.geometry; // original geometry
		var oe = og.vertices[1]; // original end

		var geom = new THREE.Geometry();
		var start = new THREE.Vector3( oe.x, oe.y, oe.z );
		var end = new THREE.Vector3( oe.x, oe.y, oe.z );
		geom.vertices.push(start, end);
		g.line = new THREE.Line(geom, mat);
		_scene.add( g.line );
	};

	// increase line width 
	this.incWidth = function() {
		_mats["go_blocked"].linewidth++;
		_mats["go_sleep"].linewidth++;
		_mats["go_normal"].linewidth++;
		_mats["go_link"].linewidth++;
		_mats["send_arrow"].linewidth++;
	};
	// increase line width 
	this.decWidth = function() {
		_mats["go_blocked"].linewidth--;
		_mats["go_sleep"].linewidth--;
		_mats["go_normal"].linewidth--;
		_mats["go_link"].linewidth--;
		_mats["send_arrow"].linewidth--;
	};

	_widths = {
		"default": [2, 2, 2, 2, 2],
		"mode1": [10, 1, 1, 1, 1],
		"mode2": [1, 10, 1, 1, 1],
		"mode3": [1, 1, 10, 1, 1],
		"mode4": [1, 1, 1, 1, 10],
	};

	this.highlight = function(mode) {
		_mats["go_blocked"].linewidth = _widths[mode][0];
		_mats["go_sleep"].linewidth =  _widths[mode][1];
		_mats["go_normal"].linewidth =  _widths[mode][2];
		_mats["go_link"].linewidth =  _widths[mode][3];
		_mats["send_arrow"].linewidth =  _widths[mode][4];
	};
};
