// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC721 {
    function transferFrom(address from, address to, uint256 tokenid) external;
}

contract DutchAuction {
    uint256 private constant DURATION = 7 days;

    IERC721 public immutable nft;
    uint256 public immutable nftId;

    address payable public immutable seller;
    uint256 public immutable startPrice;
    uint256 public immutable startTime;
    uint256 public immutable expiresAt;
    uint256 public immutable discountRate;

    constructor(uint256 _startPrice, uint256 _discountRate, address _nft, uint256 _nftId) {
        require(_startPrice >= _discountRate * DURATION, "invalid price and discount");

        seller = payable(msg.sender);
        startPrice = _startPrice;
        startTime = block.timestamp;
        expiresAt = block.timestamp + DURATION;
        discountRate = _discountRate;
        nft = IERC721(_nft);
        nftId = _nftId;
    }

    function getPrice() public view returns(uint256) {
        uint256 timeElapsed = block.timestamp - startTime;
        uint256 discount = timeElapsed * discountRate;
        return startPrice - discount;
    }

    function buy() external payable {
        require(block.timestamp < expiresAt, "auction expired");

        uint price = getPrice();
        require(msg.value >= price, "not enough value");

        nft.transferFrom(seller, msg.sender, nftId);

        uint refund = msg.value - price;
        if (refund > 0) {
            payable(msg.sender).transfer(refund);
        }

        selfdestruct(seller);
    }
}

contract EnglishAuction {
    uint256 private constant DURATION = 7 days;

    IERC721 public immutable nft;
    uint256 public immutable nftId;

    address payable public immutable seller;
    uint256 public endAt;
    bool public started;
    bool public ended;

    address public highestBidder;
    uint256 public highestBid;
    mapping(address => uint256) public bids;

    event Start();
    event Bid(address indexed sender, uint256 amount);
    event Withdraw(address indexed sender, uint256 amount);
    event End(address highestBidder, uint256 highestBid);

    constructor(address _nft, uint256 _nftId, uint256 _startingBid) {
        nft = IERC721(_nft);
        nftId = _nftId;
        seller = payable(msg.sender);
        highestBid = _startingBid;
    }

    function start() external {
        require(msg.sender == seller, "not seller");
        require(!started, "started");

        started = true;
        endAt = block.timestamp + DURATION;
        nft.transferFrom(seller, address(this), nftId);

        emit Start();
    }

    function bid() external payable {
        require(started, "not started");
        require(block.timestamp < endAt, "ended");
        require(msg.value > highestBid, "not enough value");

        if (highestBidder != address(0)) {
            bids[highestBidder] += highestBid;
        }

        highestBid = msg.value;
        highestBidder = msg.sender;

        emit Bid(msg.sender, msg.value);
    }

    function withdraw() external {
        uint256 bal = bids[msg.sender];
        bids[msg.sender] = 0;
        payable(msg.sender).transfer(bal);
        
        emit Withdraw(msg.sender, bal);
    }

    function end() external {
        require(started, "not started");
        require(!ended, "ended");
        require(block.timestamp >= endAt, "not ended");

        ended = true;

        if (highestBidder != address(0)) {
            nft.transferFrom(address(this), highestBidder, nftId);
            seller.transfer(highestBid);
        } else {
            nft.transferFrom(address(this), seller, nftId);
        }

        emit End(highestBidder, highestBid);
    }
}

