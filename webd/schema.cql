CREATE TABLE metrics (
    name  text,
    time  bigint,
    major int,
    minor int,
    rcp   int,
    rmg   int,
    rsc   int,
    rms   int,
    wcp   int,
    wmg   int,
    wsc   int,
    wms   int,
    iop   int,
    ioms  int,
    ioq   int,
    PRIMARY KEY(name,time)
);
