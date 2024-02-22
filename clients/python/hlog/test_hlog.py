# pyright: reportMissingImports=false

import pytest
from .main import Config, Hlog

@pytest.fixture
def configs() -> Config:
    return Config(
        client_id="test",
        kafka_server="0.0.0.0:65007",
        channel_id="subnet"
    )


@pytest.fixture
def hlog_test(configs) -> Hlog:
    return Hlog(configs=configs)


def test_publish(hlog_test) -> None:
    r = hlog_test._publish(
        message="hello world",
        data={"foo": "bar"}
    )
    assert r == True


def test_debug(hlog_test) -> None:
    assert hlog_test.debug("hello world")== True


def test_info(hlog_test) -> None:
    assert hlog_test.info("hello world") == True


def test_warn(hlog_test) -> None:
    assert hlog_test.warn("hello world") == True


def test_error(hlog_test) -> None:
    assert hlog_test.error("hello world") == True


def test_fatal(hlog_test) -> None:
    assert hlog_test.fatal("hello world") == True
