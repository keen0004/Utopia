// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract CrowdFund {
    struct Campaign {
        address creator;
        uint256 goal;
        uint256 pledged;
        uint256 startAt;
        uint256 endAt;
        bool claimed;
    }

    IERC20 public immutable token;
    uint256 public count;
    mapping(uint256 => Campaign) public campaigns;
    mapping(uint256 => mapping(address => uint256)) public pledgedAmount;

    event Launch(uint256 count, address caller, uint256 _goal, uint256 _startAt, uint256 _endAt);
    event Cancel(uint256 _id);
    event Pledge(uint256 indexed _id, address indexed caller, uint256 _amount);
    event Unpledge(uint256 indexed _id, address indexed caller, uint256 _amount);
    event Claim(uint256 _id);
    event Refund(uint256 indexed _id, address indexed caller, uint256 _amount);

    constructor(address _token) {
        token = IERC20(_token);
    }

    function launch(uint256 _goal, uint256 _startAt, uint256 _endAt) external {
        require(_startAt >= block.timestamp, "invalid start time");
        require(_endAt >= block.timestamp, "invalid end time");
        require(_endAt <= block.timestamp + 90 days, "over max end time");

        count += 1;
        campaigns[count] = Campaign({
            creator: msg.sender, 
            goal: _goal, 
            pledged: 0, 
            startAt: _startAt, 
            endAt: _endAt, 
            claimed: false
        });

        emit Launch(count, msg.sender, _goal, _startAt, _endAt);
    }

    function cancel(uint256 _id) external {
        Campaign storage campaign = campaigns[_id];
        
        require(msg.sender == campaign.creator, "not creator");
        require(block.timestamp < campaign.startAt, "started");
        
        delete campaigns[_id];
        emit Cancel(_id);
    }

    function pledge(uint256 _id, uint256 _amount) external {
        Campaign storage campaign = campaigns[_id];

        require(block.timestamp >= campaign.startAt, "not started");
        require(block.timestamp <= campaign.endAt, "ended");

        campaign.pledged += _amount;
        pledgedAmount[_id][msg.sender] += _amount;
        token.transferFrom(msg.sender, address(this), _amount);

        emit Pledge(_id, msg.sender, _amount);
    }

    function unpledge(uint256 _id, uint256 _amount) external {
        Campaign storage campaign = campaigns[_id];

        require(block.timestamp <= campaign.endAt, "ended");
        require(pledgedAmount[_id][msg.sender] >= _amount, "invalid amount");

        campaign.pledged -= _amount;
        pledgedAmount[_id][msg.sender] -= _amount;
        token.transferFrom(address(this), msg.sender, _amount);

        emit Unpledge(_id, msg.sender, _amount);
    }

    function claim(uint256 _id) external {
        Campaign storage campaign = campaigns[_id];

        require(msg.sender == campaign.creator, "not creator");
        require(block.timestamp > campaign.endAt, "not ended");
        require(campaign.pledged >= campaign.goal, "not enough pledged");
        require(!campaign.claimed, "claimed");

        campaign.claimed = true;
        token.transferFrom(address(this), msg.sender, campaign.pledged);

        emit Claim(_id);
    }

    function refund(uint256 _id) external {
        Campaign storage campaign = campaigns[_id];

        require(block.timestamp > campaign.endAt, "not ended");
        require(campaign.pledged < campaign.goal, "not enough pledged");

        uint256 bal = pledgedAmount[_id][msg.sender];
        pledgedAmount[_id][msg.sender] = 0;
        token.transferFrom(address(this), msg.sender, bal);

        emit Refund(_id, msg.sender, bal);
    }
}

