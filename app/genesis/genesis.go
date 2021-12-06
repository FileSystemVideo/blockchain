package genesis

const GenesisJson = `{
  "genesis_time": "2021-10-21T16:10:29.9874294Z",
  "chain_id": "fsv20211021",
  "initial_height": "1",
  "consensus_params": {
    "block": {
      "max_bytes": "22020096",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {}
  },
  "app_hash": "",
  "app_state": {
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      },
      "accounts": [
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "fsv1uvktzxpl8u4rxwhd757ury7xwprfvjplnd6cp5",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "fsv16knzs948zchx8dlaxl5tey5hs0hxr2rzj2uvfa",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "fsv18ernyc0p8gdc7tuv2qjrn4vzr67j33uenp5p7w",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        }
      ]
    },
    "bank": {
      "params": {
        "send_enabled": [],
        "default_send_enabled": true
      },
      "balances": [
        {
          "address": "fsv18ernyc0p8gdc7tuv2qjrn4vzr67j33uenp5p7w",
          "coins": [
            {
              "denom": "fsv",
              "amount": "2078999996088000"
            }
          ]
        },
        {
          "address": "fsv16knzs948zchx8dlaxl5tey5hs0hxr2rzj2uvfa",
          "coins": [
            {
              "denom": "fsv",
              "amount": "3500000000000"
            }
          ]
        },
        {
          "address": "fsv1uvktzxpl8u4rxwhd757ury7xwprfvjplnd6cp5",
          "coins": [
            {
              "denom": "fsv",
              "amount": "10500000000000"
            }
          ]
        }
      ],
      "supply": [],
      "denom_metadata": []
    },
    "capability": {
      "index": "1",
      "owners": []
    },
    "copyright": {
      "accountSpace": [],
      "accountInvite": [],
      "deflationInfor": {
        "MinerTotalAmount": "",
        "HasMinerAmount": "",
        "RemainMinerAmount": "",
        "DayMinerAmount": "",
        "DayMinerRemain": "",
        "SpaceMinerAmount": "",
        "SpaceMinerBonus": "",
        "DeflationSpaceTotal": ""
      },
      "inviteRecords": []
    },
    "crisis": {
      "constant_fee": {
        "denom": "fsv",
        "amount": "1000"
      }
    },
    "distribution": {
      "params": {
        "community_tax": "0.020000000000000000",
        "base_proposer_reward": "0.010000000000000000",
        "bonus_proposer_reward": "0.040000000000000000",
        "withdraw_addr_enabled": true
      },
      "fee_pool": {
        "community_pool": []
      },
      "delegator_withdraw_infos": [],
      "previous_proposer": "",
      "outstanding_rewards": [],
      "validator_accumulated_commissions": [],
      "validator_historical_rewards": [],
      "validator_current_rewards": [],
      "delegator_starting_infos": [],
      "validator_slash_events": []
    },
    "evidence": {
      "evidence": []
    },
    "genutil": {
      "gen_txs": [
        {
          "body": {
            "messages": [
              {
                "@type": "/cosmos.staking.v1beta1.MsgCreateValidator",
                "description": {
                  "moniker": "Honeycomb tissue",
                  "identity": "",
                  "website": "",
                  "security_contact": "",
                  "details": ""
                },
                "commission": {
                  "rate": "0.990000000000000000",
                  "max_rate": "0.990000000000000000",
                  "max_change_rate": "0.000100000000000000"
                },
                "min_self_delegation": "6000000000",
                "delegator_address": "fsv1uvktzxpl8u4rxwhd757ury7xwprfvjplnd6cp5",
                "validator_address": "fsvvaloper1uvktzxpl8u4rxwhd757ury7xwprfvjpljva7rs",
                "pubkey": {
                  "@type": "/cosmos.crypto.ed25519.PubKey",
                  "key": "lk0aJyuCoKnX0hZY1IUvFQpFkh96yfB4efZ9Pi5ZdOU="
                },
                "value": {
                  "denom": "fsv",
                  "amount": "10000000000"
                }
              }
            ],
            "memo": "061412339c8237fe8d80108a59a56ceaa018ceaa@192.168.0.38:26656",
            "timeout_height": "0",
            "extension_options": [],
            "non_critical_extension_options": []
          },
          "auth_info": {
            "signer_infos": [
              {
                "public_key": {
                  "@type": "/cosmos.crypto.secp256k1.PubKey",
                  "key": "A8gc16YOF4k9eYyZQ7Dypo1u3kf6jy1PPz7mrw1MHeIN"
                },
                "mode_info": {
                  "single": {
                    "mode": "SIGN_MODE_DIRECT"
                  }
                },
                "sequence": "0"
              }
            ],
            "fee": {
              "amount": [],
              "gas_limit": "2000000",
              "payer": "",
              "granter": ""
            }
          },
          "signatures": [
            "DSD++Aj5SMDWv3o5VKYe5bwDaqbtnpPrp/qnRIDeYA0ond1P5A5p+C6varBBOwYqb/DOA6OKgbm6ZVpevo/EFw=="
          ]
        }
      ]
    },
    "gov": {
      "starting_proposal_id": "1",
      "deposits": [],
      "votes": [],
      "proposals": [],
      "deposit_params": {
        "min_deposit": [
          {
            "denom": "fsv",
            "amount": "10000000"
          }
        ],
        "max_deposit_period": "172800s"
      },
      "voting_params": {
        "voting_period": "172800s"
      },
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000",
        "veto_threshold": "0.334000000000000000"
      }
    },
    "ibc": {
      "client_genesis": {
        "clients": [],
        "clients_consensus": [],
        "clients_metadata": [],
        "params": {
          "allowed_clients": [
            "06-solomachine",
            "07-tendermint"
          ]
        },
        "create_localhost": false,
        "next_client_sequence": "0"
      },
      "connection_genesis": {
        "connections": [],
        "client_connection_paths": [],
        "next_connection_sequence": "0"
      },
      "channel_genesis": {
        "channels": [],
        "acknowledgements": [],
        "commitments": [],
        "receipts": [],
        "send_sequences": [],
        "recv_sequences": [],
        "ack_sequences": [],
        "next_channel_sequence": "0"
      }
    },
    "params": null,
    "slashing": {
      "params": {
        "signed_blocks_window": "600",
        "min_signed_per_window": "0.500000000000000000",
        "downtime_jail_duration": "1800s",
        "slash_fraction_double_sign": "0.050000000000000000",
        "slash_fraction_downtime": "0.001000000000000000"
      },
      "signing_infos": [],
      "missed_blocks": []
    },
    "staking": {
      "params": {
        "unbonding_time": "300s",
        "max_validators": 51,
        "max_entries": 7,
        "historical_entries": 10000,
        "bond_denom": "fsv"
      },
      "last_total_power": "0",
      "last_validator_powers": [],
      "validators": [],
      "delegations": [],
      "unbonding_delegations": [],
      "redelegations": [],
      "exported": false
    },
    "transfer": {
      "port_id": "transfer",
      "denom_traces": [],
      "params": {
        "send_enabled": true,
        "receive_enabled": true
      }
    },
    "upgrade": {},
    "vesting": {}
  }
}`
