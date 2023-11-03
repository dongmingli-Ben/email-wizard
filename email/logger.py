"""
A minimal Python thread-safe, exception capturing logging module.

Usage:

>>> #### somewhere in your main function, before you try to log anything
>>> import logger
>>> logger.logger_init(log_dir='logs', name='app', when='D', 
                   backupCount=7, level='INFO')   # it is critical to initialize the logger before using it

>>> #### in some other file where you want to log
>>> import logger
>>> logger.critical("this is a critical message")
>>> logger.fatal("this is a fatal message")
>>> logger.error("this is a error message")
>>> logger.warn("this is a warn message")
>>> logger.warning("this is a warning message")
>>> logger.info("this is a info message")
>>> logger.debug("this is a debug message")

critical/fatal/error will also print the stack trace of the exception:
>>> def error_func():
>>>     try:
>>>         1 / 0
>>>     except ZeroDivisionError as e:
>>>         fatal("error with zero division")
>>>     except:
>>>         error("error!")
>>>     error("calling error outside exception")
>>> 
>>> logger_init()
>>> error_func()
2023-11-04 03:40:12,587 - app - CRITICAL - error with zero division
Traceback (most recent call last):
  File "test.py", line 74, in error_func
    1 / 0
ZeroDivisionError: division by zero
2023-11-04 03:40:12,588 - app - ERROR - calling error outside exception
NoneType: None
"""

import logging
import logging.handlers as handlers
import os

LOGGING_LEVEL_MAP = {
    'DEBUG': logging.DEBUG,
    'INFO': logging.INFO,
    'WARN': logging.WARN,
    'ERROR': logging.ERROR,
    'CRITICAL': logging.CRITICAL
}

logger = None


def logger_init(log_dir='log', name='app', when='D', backupCount=7, level='INFO'):
    """
    This will initialize the logger which logs both to the terminal and a file.
    The parameters are similar to Python's logging module.

    Note: logger initialization is not thread-safe. If you would like to use logger in 
    multiple threads, initialize the logger outside the threads.
    """
    global logger
    if logger is not None:
        logger.warning("logger received re-init call, ignoring ...")
        return
    logger = logging.getLogger(name)
    if level not in LOGGING_LEVEL_MAP:
        logger.warning(
            f"logging level {level} not acceptable, setting level to INFO ...")
        level = 'INFO'
    logger.setLevel(LOGGING_LEVEL_MAP[level])
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    if not os.path.exists(log_dir):
        os.makedirs(log_dir)
    file_log_handler = handlers.TimedRotatingFileHandler(
        os.path.join(log_dir, f'{name}.log'), when=when, interval=1, backupCount=backupCount)
    file_log_handler.setLevel(LOGGING_LEVEL_MAP[level])
    # Here we set our logHandler's formatter
    file_log_handler.setFormatter(formatter)
    logger.addHandler(file_log_handler)
    terminal_log = logging.StreamHandler()
    terminal_log.setLevel(LOGGING_LEVEL_MAP[level])
    terminal_log.setFormatter(formatter)
    logger.addHandler(terminal_log)


def info(*args, **kwargs):
    logger.info(*args, **kwargs)


def debug(*args, **kwargs):
    logger.debug(*args, **kwargs)


def warning(*args, **kwargs):
    logger.warning(*args, **kwargs)


def error(*args, **kwargs):
    kwargs['exc_info'] = True
    logger.error(*args, **kwargs)


def critical(*args, **kwargs):
    kwargs['exc_info'] = True
    logger.critical(*args, **kwargs)


warn = warning
fatal = critical
