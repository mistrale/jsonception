### begin wb_initialization
### end wb_initialization
### begin wb_common
### end wb_common
### begin wb_destroy
### end wb_destroy	
	
	
	
def wb_export(f):
    import json, inspect
    params = json.loads(f.__doc__)
    if "in" not in params and "out" not in params:
        raise ValueError("malformed parameter")
    try:
        #remove type validation
        pass
        #if "value" not in params['out']:
        #    raise ValueError("out should contain value")
        #for k, v in params['in'].iteritems():
        #    if (not isinstance(v, basestring)) and (not isinstance(v, int)) \
        #            and (not isinstance(v, dict)) and (not isinstance(v, list)) \
        #            and (not isinstance(v, bool)) and (not isinstance(v, float)):
        #        raise ValueError("not string")
    except AttributeError as e:
        raise ValueError(e)
    except TypeError as e:
        raise ValueError(e)
    __wb_export = True
    def json_wrapper(input):
        input_d = json.loads(input)
        if 'd' in inspect.getargspec(f).args:
            return_value = f(input, d=input_d)
        else:
            return_value = f(input)
        if isinstance(return_value, str):
            return return_value
        else:
            return json.dumps(return_value)
    json_wrapper.__wb_export = __wb_export
    json_wrapper.__doc__ = getattr(f, '__doc__', None)

    return json_wrapper

@wb_export
def unnamed_block1(jsonstr, d={}):
    # Definition of Input and Output Parameters (name and type)
    """{"in":{},"out":{}}"""
    import json
    # You code here
    ret_dict = {}
    ret_dict["returncode"] = 0 # set to 1 to output on the FAILURE output
    print json.dumps(d)
    return ret_dict