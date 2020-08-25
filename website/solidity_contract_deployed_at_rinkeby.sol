pragma solidity >= 0.4.22 < 0.7.0;


contract Storage {

    struct Action{
        uint256 timestamp;
        string action_happening;
    }

    mapping(string => Action[]) assets;

    function store(string memory _key, string memory _value) public {
        assets[_key].push(Action(now,_value));
    }

    function length(string memory _key) public view returns (uint256){
       return assets[_key].length;  
    }

    function retrieve(string memory _key, uint256 index) public view returns (uint256, string memory){
        return (assets[_key][index].timestamp, assets[_key][index].action_happening);    
    }
}
