from flask import Flask, request, jsonify, Blueprint, Response
from flask_cors import cross_origin
import json
import ui_api_entity_config
import flask_utils

blueprint_route_extract_ui_api = Blueprint('blueprint_route_extract_ui_api', __name__)

@blueprint_route_extract_ui_api.route('/extract_ui_api_entity', methods=['POST'])
@cross_origin(origin='*')
def extract_ui_api_entity():
    use_cache = flask_utils.get_arg_bool('use_cache', False)
    tenant_key = flask_utils.get_arg('tenant_key', '0')
    entity_id = flask_utils.get_arg('entity_id')

    done = ui_api_entity_config.get_entity(
        tenant_key, entity_id, use_cache, cache_only=False)

    response = Response(
        response=json.dumps(done),
        status=200,
        mimetype='application/json'
    )

    return response
